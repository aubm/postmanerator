package postman

import (
	"reflect"
	"testing"
)

func TestExtractStructuresDefinition(t *testing.T) {
	// Given
	builder := &CollectionBuilder{}
	builder.Parsers = append(builder.Parsers, &CollectionV210Parser{})
	expectedStructures := []StructureDefinition{
		{Name: "Cat", Description: "A great animal", Fields: []StructureFieldDefinition{
			{Name: "id", Description: "A unique identifier for the cat", Type: "int"},
			{Name: "color", Description: "The color of the cat", Type: "string"},
			{Name: "name", Description: "The name of the cat", Type: "string"},
		}},
		{Name: "Dog", Description: "A greater animal", Fields: []StructureFieldDefinition{
			{Name: "id", Description: "A unique identifier for the dog", Type: "int"},
			{Name: "color", Description: "The color of the dog", Type: "string"},
			{Name: "name", Description: "The name of the dog", Type: "string"},
		}},
	}

	// When
	col, err := builder.FromFile("tests_data/collection-01.json", BuilderOptions{})

	// Then
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if reflect.DeepEqual(col.Structures, expectedStructures) == false {
		t.Errorf("Collection structures definition were not properly extracted, expected %v, got %v",
			expectedStructures, col.Structures)
	}
}
