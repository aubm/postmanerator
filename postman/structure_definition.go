package postman

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
)

type StructureDefinition struct {
	Name        string
	Description string
	Fields      []StructureFieldDefinition
}

type StructureFieldDefinition struct {
	Name        string
	Description string
	Type        string
}

func (col *Collection) ExtractStructuresDefinition() {
	vm := otto.New()
	var codeFragments []string
	for _, req := range col.Requests {
		codeFragments = append(codeFragments, extractCodeFragments(req.Tests)...)
	}
	vm.Set("APIStructures", struct{}{})
	for _, frag := range codeFragments {
		frag = frag + `
if (!!populateNewAPIStructures && typeof(populateNewAPIStructures) === 'function') {
    populateNewAPIStructures();
}`
		vm.Run(frag)
	}
	if value, err := vm.Get("APIStructures"); err == nil {
		if apiStructures := value.Object(); apiStructures != nil {
			for _, key := range apiStructures.Keys() {
				if structureDef, err := apiStructures.Get(key); err == nil {
					if structure, err := getStructureDefinition(structureDef); err == nil {
						col.Structures = append(col.Structures, structure)
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

func extractCodeFragments(input string) []string {
	var codeFragments []string
	var fragment string
	validID := regexp.MustCompile(`\/\*\[\[(start|end) postmanerator\]\]\*\/`)
	for validID.MatchString(input) {
		fragment, input = nextCodeFragment(input)
		codeFragments = append(codeFragments, fragment)
	}

	return codeFragments
}

func nextCodeFragment(input string) (string, string) {
	startTag := "/*[[start postmanerator]]*/"
	endTag := "/*[[end postmanerator]]*/"
	parts := strings.Split(input, startTag)
	parts = parts[1:]
	input = strings.Join(parts, startTag)
	parts = strings.Split(input, endTag)
	return strings.Trim(parts[0], "\n"), strings.Join(parts[1:], endTag)
}

func getStructureDefinition(srcVal otto.Value) (StructureDefinition, error) {
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
