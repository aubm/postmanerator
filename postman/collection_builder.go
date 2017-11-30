package postman

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
)

var (
	ErrRequestNotFound  = errors.New("request not found")
	ErrAllParsersFailed = errors.New("failed to parse collection file, all parsers failed")
)

type CollectionBuilder struct {
	Parsers []interface {
		CanParse(contents []byte) bool
		Parse(contents []byte, options BuilderOptions) (Collection, error)
	}
}

func (c *CollectionBuilder) FromFile(file string, options BuilderOptions) (col Collection, err error) {
	b, err := c.readCollectionFile(file, options)
	if err != nil {
		return col, err
	}

	col, err = c.parseCollection(b, options)
	if err != nil {
		return col, err
	}

	c.extractStructuresDefinition(&col)
	return col, nil
}

func (c *CollectionBuilder) readCollectionFile(file string, options BuilderOptions) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if options.EnvironmentVariables != nil {
		for k, v := range options.EnvironmentVariables {
			b = bytes.Replace(b, []byte(fmt.Sprintf("{{%v}}", k)), []byte(v), -1)
		}
	}

	return b, nil
}

func (c *CollectionBuilder) parseCollection(contents []byte, options BuilderOptions) (Collection, error) {
	for _, p := range c.Parsers {
		if p.CanParse(contents) {
			return p.Parse(contents, options)
		}
	}
	return Collection{}, ErrAllParsersFailed
}

func (c *CollectionBuilder) extractStructuresDefinition(col *Collection) {
	structureDefinitions := make([]StructureDefinition, 0)

	tests := c.extractCollectionTests(col)

	codeFragments := make([]string, 0)
	for _, t := range tests {
		codeFragments = append(codeFragments, c.extractCodeFragments(t)...)
	}

	vm := otto.New()
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

	col.Structures = structureDefinitions
}

func (c *CollectionBuilder) extractCollectionTests(col *Collection) []string {
	f := Folder{Folders: col.Folders, Requests: col.Requests}
	return c.extractFolderTests(f)
}

func (c *CollectionBuilder) extractFolderTests(folder Folder) []string {
	tests := make([]string, 0)
	for _, req := range folder.Requests {
		tests = append(tests, req.Tests)
	}
	for _, f := range folder.Folders {
		tests = append(tests, c.extractFolderTests(f)...)
	}
	return tests
}

func (c *CollectionBuilder) extractCodeFragments(input string) []string {
	var codeFragments []string
	var fragment string
	validID := regexp.MustCompile(`/\*\[\[(start|end) postmanerator]]\*/`)
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
		return structDef, errors.New("value is not an object")
	}

	srcObj := srcVal.Object()

	// Get structure name
	nameVal, err := srcObj.Get("name")
	if err != nil {
		return structDef, errors.New("structures must have a name")
	}
	structDef.Name = nameVal.String()

	// Get structure description
	descVal, err := srcObj.Get("description")
	if err != nil {
		return structDef, errors.New("structures must have a description")
	}
	structDef.Description = descVal.String()

	// Get structure fields
	fieldsVal, err := srcObj.Get("fields")
	if err != nil {
		return structDef, errors.New("structures must define fields")
	}
	fields, err := fieldsVal.Export()
	if err != nil {
		return structDef, errors.New("failed to convert javascript fields attribute into a valid go type")
	}
	fieldsSlice, ok := fields.([]map[string]interface{})
	if !ok {
		return structDef, errors.New("fields attribute must be an array of objects")
	}
	for _, fieldDefMap := range fieldsSlice {
		// Get field name
		var fieldDef StructureFieldDefinition
		if fieldName, ok := fieldDefMap["name"].(string); ok == true {
			fieldDef.Name = fieldName
		} else {
			return structDef, errors.New("structure fields must have a name")
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
