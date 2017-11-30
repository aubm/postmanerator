package postman

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CollectionV1Parser struct{}

func (p *CollectionV1Parser) CanParse(contents []byte) bool {
	return true
}

func (p *CollectionV1Parser) Parse(contents []byte, options BuilderOptions) (Collection, error) {
	src := collectionV1{}
	if err := json.Unmarshal(contents, &src); err != nil {
		return Collection{}, err
	}
	return p.buildCollectionFromV1(src, options)
}

func (p *CollectionV1Parser) buildCollectionFromV1(src collectionV1, options BuilderOptions) (Collection, error) {
	collection := Collection{
		Name:        src.Name,
		Description: src.Description,
		Requests:    make([]Request, 0),
		Folders:     make([]Folder, 0),
		Structures:  make([]StructureDefinition, 0),
	}

	for _, requestID := range src.Order {
		req, err := p.buildRequest(src, requestID, options)
		if err != nil {
			return collection, fmt.Errorf("failed to build request %v: %v", requestID, err)
		}
		collection.Requests = append(collection.Requests, req)
	}

	for _, folder := range src.Folders {
		newFolder := Folder{
			ID:          folder.ID,
			Name:        folder.Name,
			Description: folder.Description,
			Requests:    make([]Request, 0),
		}
		for _, requestID := range folder.Order {
			req, err := p.buildRequest(src, requestID, options)
			if err != nil {
				return collection, fmt.Errorf("failed to build request %v: %v", requestID, err)
			}
			newFolder.Requests = append(newFolder.Requests, req)
		}
		collection.Folders = append(collection.Folders, newFolder)
	}

	return collection, nil
}

func (p *CollectionV1Parser) buildRequest(colV1 collectionV1, requestID string, options BuilderOptions) (Request, error) {
	request := Request{}

	v1, ok := p.findRequestByID(colV1, requestID)
	if !ok {
		return request, ErrRequestNotFound
	}

	request.ID = v1.ID
	request.Name = v1.Name
	request.Description = v1.Description
	request.Method = v1.Method
	request.URL = v1.URL
	request.Tests = v1.Tests
	request.PayloadType = v1.DataMode
	request.PayloadRaw = v1.RawModeData
	request.PayloadParams = p.buildRequestPayloadParams(v1)
	request.PathVariables = p.buildRequestPathVariables(v1)
	request.Headers = p.buildRequestHeaders(v1, options)
	request.Responses = p.buildRequestResponses(v1, options)

	return request, nil
}

func (p *CollectionV1Parser) findRequestByID(v1 collectionV1, id string) (requestV1, bool) {
	for _, req := range v1.Requests {
		if req.ID == id {
			return req, true
		}
	}
	return requestV1{}, false
}

func (p *CollectionV1Parser) buildRequestPayloadParams(v1 requestV1) []KeyValuePair {
	if v1.Data == nil {
		return nil
	}

	payloadParams := make([]KeyValuePair, 0)
	for _, d := range v1.Data {
		payloadParams = append(payloadParams, KeyValuePair{
			Name:  d.Key,
			Key:   d.Key,
			Value: d.Value,
		})
	}

	return payloadParams
}

func (p *CollectionV1Parser) buildRequestPathVariables(v1 requestV1) []KeyValuePair {
	if v1.PathVariables == nil {
		return nil
	}

	pathVariables := make([]KeyValuePair, 0)
	for name, value := range v1.PathVariables {
		pathVariables = append(pathVariables, KeyValuePair{
			Name:  name,
			Key:   name,
			Value: value,
		})
	}

	return pathVariables
}

func (p *CollectionV1Parser) buildRequestHeaders(v1 requestV1, options BuilderOptions) []KeyValuePair {
	headers := make([]KeyValuePair, 0)
	rawHeaderList := strings.Split(v1.RawHeaders, "\n")
	for _, rawHeader := range rawHeaderList {
		parts := strings.Split(rawHeader, ": ")
		if len(parts) != 2 || p.contains(options.IgnoredRequestHeaders, parts[0]) {
			continue
		}
		headers = append(headers, KeyValuePair{
			Name:  parts[0],
			Key:   parts[0],
			Value: parts[1],
		})
	}
	return headers
}

func (p *CollectionV1Parser) buildRequestResponses(v1 requestV1, options BuilderOptions) []Response {
	responses := make([]Response, 0)
	for _, res := range v1.Responses {
		responses = append(responses, Response{
			ID:         res.ID,
			Name:       res.Name,
			Status:     res.Status,
			StatusCode: res.ResponseCode.Code,
			Body:       res.Text,
			Headers:    p.buildResponseHeaders(res, options),
		})
	}

	return responses
}

func (p *CollectionV1Parser) buildResponseHeaders(v1 responseV1, options BuilderOptions) []KeyValuePair {
	headers := make([]KeyValuePair, 0)
	for _, header := range v1.Headers {
		if p.contains(options.IgnoredResponseHeaders, header.Name) {
			continue
		}
		headers = append(headers, KeyValuePair{
			Name:        header.Name,
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
		})
	}
	return headers
}

func (p *CollectionV1Parser) contains(target []string, symbol string) bool {
	for _, element := range target {
		if element == symbol {
			return true
		}
	}
	return false
}
