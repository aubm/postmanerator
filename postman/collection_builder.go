package postman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
)

var (
	ErrRequestNotFound = errors.New("request not found")
	ErrFolderNotFound  = errors.New("folder not found")
)

type CollectionBuilder struct{}

func (c *CollectionBuilder) FromFile(file string, options BuilderOptions) (Collection, error) {
	parsedCollection, err := c.parseCollectionFile(file, options)
	if err != nil {
		return Collection{}, err
	}

	return c.buildCollectionFromV1(parsedCollection, options)
}

func (c *CollectionBuilder) parseCollectionFile(file string, options BuilderOptions) (collectionV1, error) {
	collection := collectionV1{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return collection, err
	}

	if options.EnvironmentVariables != nil {
		for k, v := range options.EnvironmentVariables {
			b = bytes.Replace(b, []byte(fmt.Sprintf("{{%v}}", k)), []byte(v), -1)
		}
	}

	err = json.Unmarshal(b, &collection)
	if err != nil {
		return collection, err
	}

	return collection, nil
}

func (c *CollectionBuilder) buildCollectionFromV1(src collectionV1, options BuilderOptions) (Collection, error) {
	collection := Collection{
		Name:        src.Name,
		Description: src.Description,
		Requests:    make([]Request, 0),
		Folders:     make([]Folder, 0),
		Structures:  make([]StructureDefinition, 0),
	}

	for _, requestID := range src.Order {
		req, err := c.buildRequest(src, requestID, options)
		if err != nil {
			return collection, fmt.Errorf("failed to build request %v: %v", requestID, err)
		}
		collection.Requests = append(collection.Requests, req)
	}

	if len(src.FoldersOrder) > 0 {
		for _, folderID := range src.FoldersOrder {

			folder, err := c.buildFolder(src, folderID, options)
			if err != nil {
				return collection, fmt.Errorf("failed to build folder %v: %v", folderID, err)
			}
			collection.Folders = append(collection.Folders, folder)
		}
	} else {
		for _, folder := range src.Folders {
			newFolder := Folder{
				ID:          folder.ID,
				Name:        folder.Name,
				Description: folder.Description,
				Requests:    make([]Request, 0),
			}
			for _, requestID := range folder.Order {
				req, err := c.buildRequest(src, requestID, options)
				if err != nil {
					return collection, fmt.Errorf("failed to build request %v: %v", requestID, err)
				}
				newFolder.Requests = append(newFolder.Requests, req)
			}
			collection.Folders = append(collection.Folders, newFolder)
		}
	}

	collection.Structures = c.extractStructuresDefinition(src)

	return collection, nil
}

func (c *CollectionBuilder) buildFolder(colV1 collectionV1, folderID string, options BuilderOptions) (Folder, error) {
	folder := Folder{}

	fol, ok := c.findFolderByID(colV1, folderID)
	if !ok {
		return folder, ErrFolderNotFound
	}

	folder.ID = fol.ID
	folder.Name = fol.Name
	folder.Description = fol.Description
	folder.Requests = make([]Request, 0)

	for _, requestID := range fol.Order {
		req, err := c.buildRequest(colV1, requestID, options)
		if err != nil {
			return folder, fmt.Errorf("failed to build folder %v: %v", requestID, err)
		}
		folder.Requests = append(folder.Requests, req)
	}

	return folder, nil
}

func (c *CollectionBuilder) findFolderByID(colV1 collectionV1, id string) (foldersV1, bool) {
	for _, folder := range colV1.Folders {
		if folder.ID == id {
			return folder, true
		}
	}
	return foldersV1{}, false
}

func (c *CollectionBuilder) buildRequest(colV1 collectionV1, requestID string, options BuilderOptions) (Request, error) {
	request := Request{}

	v1, ok := c.findRequestByID(colV1, requestID)
	if !ok {
		return request, ErrRequestNotFound
	}

	request.ID = v1.ID
	request.Name = v1.Name
	request.Description = v1.Description
	request.Method = v1.Method
	request.URL = v1.URL
	request.PayloadType = v1.DataMode
	request.PayloadRaw = v1.RawModeData
	request.PayloadParams = c.buildRequestPayloadParams(v1)
	request.PathVariables = c.buildRequestPathVariables(v1)
	request.Headers = c.buildRequestHeaders(v1, options)
	request.Responses = c.buildRequestResponses(v1, options)

	return request, nil
}

func (c *CollectionBuilder) findRequestByID(v1 collectionV1, id string) (requestV1, bool) {
	for _, req := range v1.Requests {
		if req.ID == id {
			return req, true
		}
	}
	return requestV1{}, false
}

func (c *CollectionBuilder) buildRequestPayloadParams(v1 requestV1) []KeyValuePair {
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

func (c *CollectionBuilder) buildRequestPathVariables(v1 requestV1) []KeyValuePair {
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

func (c *CollectionBuilder) buildRequestHeaders(v1 requestV1, options BuilderOptions) []KeyValuePair {
	headers := make([]KeyValuePair, 0)
	rawHeaderList := strings.Split(v1.RawHeaders, "\n")
	for _, rawHeader := range rawHeaderList {
		parts := strings.Split(rawHeader, ": ")
		if len(parts) != 2 || c.contains(options.IgnoredRequestHeaders, parts[0]) {
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

func (c *CollectionBuilder) buildRequestResponses(v1 requestV1, options BuilderOptions) []Response {
	responses := make([]Response, 0)
	for _, res := range v1.Responses {
		responses = append(responses, Response{
			ID:         res.ID,
			Name:       res.Name,
			Status:     res.Status,
			StatusCode: res.ResponseCode.Code,
			Body:       res.Text,
			Headers:    c.buildResponseHeaders(res, options),
		})
	}

	return responses
}

func (c *CollectionBuilder) buildResponseHeaders(v1 responseV1, options BuilderOptions) []KeyValuePair {
	headers := make([]KeyValuePair, 0)
	for _, header := range v1.Headers {
		if c.contains(options.IgnoredResponseHeaders, header.Name) {
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

func (c *CollectionBuilder) contains(target []string, symbol string) bool {
	for _, element := range target {
		if element == symbol {
			return true
		}
	}
	return false
}

func (c *CollectionBuilder) extractStructuresDefinition(v1 collectionV1) []StructureDefinition {
	structureDefinitions := make([]StructureDefinition, 0)

	vm := otto.New()
	var codeFragments []string
	for _, req := range v1.Requests {
		codeFragments = append(codeFragments, c.extractCodeFragments(req.Tests)...)
	}
	vm.Set("APIStructures", struct{}{})
	color.Set(color.FgCyan)
	for _, frag := range codeFragments {
		frag = frag + `
if (!!populateNewAPIStructures && typeof(populateNewAPIStructures) === 'function') {
    populateNewAPIStructures();
}`
		vm.Run(frag)
	}
	color.Unset()
	if value, err := vm.Get("APIStructures"); err == nil {
		if apiStructures := value.Object(); apiStructures != nil {
			for _, key := range apiStructures.Keys() {
				if structureDef, err := apiStructures.Get(key); err == nil {
					if structure, err := c.getStructureDefinition(structureDef); err == nil {
						structureDefinitions = append(structureDefinitions, structure)
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	}

	return structureDefinitions
}

func (c *CollectionBuilder) extractCodeFragments(input string) []string {
	var codeFragments []string
	var fragment string
	validID := regexp.MustCompile(`\/\*\[\[(start|end) postmanerator\]\]\*\/`)
	for validID.MatchString(input) {
		fragment, input = c.nextCodeFragment(input)
		codeFragments = append(codeFragments, fragment)
	}

	return codeFragments
}

func (c *CollectionBuilder) nextCodeFragment(input string) (string, string) {
	startTag := "/*[[start postmanerator]]*/"
	endTag := "/*[[end postmanerator]]*/"
	parts := strings.Split(input, startTag)
	parts = parts[1:]
	input = strings.Join(parts, startTag)
	parts = strings.Split(input, endTag)
	return strings.Trim(parts[0], "\n"), strings.Join(parts[1:], endTag)
}

func (c *CollectionBuilder) getStructureDefinition(srcVal otto.Value) (StructureDefinition, error) {
	var structDef StructureDefinition
	if !srcVal.IsObject() {
		return structDef, errors.New("Value is not an object")
	}

	srcObj := srcVal.Object()

	// Get structure name
	nameVal, err := srcObj.Get("name")
	if err != nil {
		return structDef, errors.New("Structures must have a name")
	}
	structDef.Name = nameVal.String()

	// Get structure description
	descVal, err := srcObj.Get("description")
	if err != nil {
		return structDef, errors.New("Structures must have a description")
	}
	structDef.Description = descVal.String()

	// Get structure fields
	fieldsVal, err := srcObj.Get("fields")
	if err != nil {
		return structDef, errors.New("Structures must define fields")
	}
	fields, err := fieldsVal.Export()
	if err != nil {
		return structDef, errors.New("Failed to convert javascript fields attribute into a valid go type")
	}
	fieldsSlice, ok := fields.([]map[string]interface{})
	if !ok {
		return structDef, errors.New("Fields attribute must be an array of objects")
	}
	for _, fieldDefMap := range fieldsSlice {
		// Get field name
		var fieldDef StructureFieldDefinition
		if fieldName, ok := fieldDefMap["name"].(string); ok == true {
			fieldDef.Name = fieldName
		} else {
			return structDef, errors.New("Structure fields must have a name")
		}

		// Get field description
		if fieldDesc, ok := fieldDefMap["description"].(string); ok == true {
			fieldDef.Description = fieldDesc
		}

		// Get field type
		if fieldType, ok := fieldDefMap["type"].(string); ok == true {
			fieldDef.Type = fieldType
		}

		structDef.Fields = append(structDef.Fields, fieldDef)
	}

	return structDef, nil
}

type BuilderOptions struct {
	IgnoredRequestHeaders  []string
	IgnoredResponseHeaders []string
	EnvironmentVariables   Environment
}
