package postman

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type CollectionV210Parser struct{}

func (p *CollectionV210Parser) CanParse(contents []byte) bool {
	return true // TODO: compare json schema value
}

func (p *CollectionV210Parser) Parse(contents []byte, options BuilderOptions) (Collection, error) {
	src := collectionV210{}
	if err := json.Unmarshal(contents, &src); err != nil {
		return Collection{}, err
	}
	return p.buildCollection(src, options)
}

func (p *CollectionV210Parser) buildCollection(src collectionV210, options BuilderOptions) (Collection, error) {
	collection := Collection{
		Name:        src.Info.Name,
		Description: src.Info.Description,
		Requests:    make([]Request, 0),
		Folders:     make([]Folder, 0),
		Structures:  make([]StructureDefinition, 0),
	}

	rootItem := Folder{}
	if err := p.computeItem(&rootItem, src.Item, options); err != nil {
		return collection, fmt.Errorf("failed to build request: %v", err)
	}

	collection.Requests = rootItem.Requests
	collection.Folders = rootItem.Folders

	return collection, nil
}

func (p *CollectionV210Parser) computeItem(parentFolder *Folder, items []collectionV210Item, options BuilderOptions) error {
	for _, item := range items {
		if item.Request == nil { // item is a folder
			folder := Folder{
				ID:          uuid.NewV4().String(),
				Description: item.Description,
				Name:        item.Name,
			}
			if err := p.computeItem(&folder, item.Item, options); err != nil {
				return err
			}
			parentFolder.Folders = append(parentFolder.Folders, folder)
		} else { // item is a request
			request := Request{
				ID:            uuid.NewV4().String(),
				Name:          item.Name,
				Description:   item.Request.Description,
				Method:        item.Request.Method,
				URL:           item.Request.Url.Raw,
				PayloadType:   item.Request.Body.Mode,
				PayloadRaw:    item.Request.Body.Raw,
				Tests:         p.parseRequestTests(item),
				QueryParams:   p.parseRequestQueryParams(item),
				PathVariables: p.parseRequestPathVariables(item),
				PayloadParams: p.parseRequestPayloadParams(item),
				Headers:       p.parseRequestHeaders(item, options),
				Responses:     p.parseRequestResponses(item, options),
			}
			parentFolder.Requests = append(parentFolder.Requests, request)
		}
	}

	return nil
}

func (p *CollectionV210Parser) parseRequestTests(item collectionV210Item) string {
	for _, event := range item.Event {
		if event.Listen == "test" {
			return strings.Join(event.Script.Exec, "\n")
		}
	}
	return ""
}

func (p *CollectionV210Parser) parseRequestPathVariables(item collectionV210Item) []KeyValuePair {
	pathVariables := make([]KeyValuePair, 0)

	for _, variable := range item.Request.Url.Variable {
		pathVariables = append(pathVariables, KeyValuePair{
			Name:        variable.Key,
			Key:         variable.Key,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	return pathVariables
}

func (p *CollectionV210Parser) parseRequestQueryParams(item collectionV210Item) []KeyValuePair {
	queryVariables := make([]KeyValuePair, 0)

	for _, variable := range item.Request.Url.Query {
		queryVariables = append(queryVariables, KeyValuePair{
			Name:        variable.Key,
			Key:         variable.Key,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	return queryVariables
}

func (p *CollectionV210Parser) parseRequestPayloadParams(item collectionV210Item) []KeyValuePair {
	payloadParams := make([]KeyValuePair, 0)

	keyValuePairCollection := make([]collectionV210KeyValuePair, 0)
	switch item.Request.Body.Mode {
	case "urlencoded":
		keyValuePairCollection = item.Request.Body.UrlEncoded
	case "formdata":
		keyValuePairCollection = item.Request.Body.FormData
	}

	for _, pair := range keyValuePairCollection {
		payloadParams = append(payloadParams, KeyValuePair{
			Name:        pair.Key,
			Key:         pair.Key,
			Value:       pair.Value,
			Description: pair.Description,
		})
	}

	return payloadParams
}

func (p *CollectionV210Parser) parseRequestHeaders(item collectionV210Item, options BuilderOptions) []KeyValuePair {
	headers := make([]KeyValuePair, 0)

	for _, header := range item.Request.Header {
		if containsString(options.IgnoredRequestHeaders, header.Key) {
			continue
		}
		headers = append(headers, KeyValuePair{
			Name:        header.Key,
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
		})
	}

	return headers
}

func (p *CollectionV210Parser) parseRequestResponses(item collectionV210Item, options BuilderOptions) []Response {
	responses := make([]Response, 0)

	for _, resp := range item.Response {
		var req collectionV210Item
		req.Request = resp.Request
		responses = append(responses, Response{
			ID:         uuid.NewV4().String(),
			Name:       resp.Name,
			Body:       resp.Body,
			Status:     resp.Status,
			StatusCode: resp.Code,
			Headers:    p.parseResponseHeaders(resp.Header, options),
			Request: Request{
				Description:   resp.Request.Description,
				Method:        resp.Request.Method,
				URL:           resp.Request.Url.Raw,
				PayloadType:   resp.Request.Body.Mode,
				PayloadRaw:    resp.Request.Body.Raw,
				QueryParams:   p.parseRequestQueryParams(req),
				PathVariables: p.parseRequestPathVariables(req),
				PayloadParams: p.parseRequestPayloadParams(req),
				Headers:       p.parseRequestHeaders(req, options),
			},
		})
	}

	return responses
}

func (p *CollectionV210Parser) parseResponseHeaders(headers []collectionV210KeyValuePair, options BuilderOptions) []KeyValuePair {
	parsedHeaders := make([]KeyValuePair, 0)

	for _, header := range headers {
		if containsString(options.IgnoredResponseHeaders, header.Key) {
			continue
		}
		parsedHeaders = append(parsedHeaders, KeyValuePair{
			Name:        header.Key,
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
		})
	}
	return parsedHeaders
}

func containsString(target []string, symbol string) bool {
	for _, element := range target {
		if element == symbol {
			return true
		}
	}
	return false
}
