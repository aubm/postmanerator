package postman

import (
	"reflect"
	"testing"
)

func TestExtractStructuresDefinition(t *testing.T) {
	// Given
	builder := &CollectionBuilder{}
	expectedStructures := []StructureDefinition{
		{Name: "Dog", Description: "A greater animal", Fields: []StructureFieldDefinition{
			{Name: "id", Description: "A unique identifier for the dog", Type: "int"},
			{Name: "color", Description: "The color of the dog", Type: "string"},
			{Name: "name", Description: "The name of the dog", Type: "string"},
		}},
		{Name: "Cat", Description: "A great animal", Fields: []StructureFieldDefinition{
			{Name: "id", Description: "A unique identifier for the cat", Type: "int"},
			{Name: "color", Description: "The color of the cat", Type: "string"},
			{Name: "name", Description: "The name of the cat", Type: "string"},
		}},
	}

	// When
	col, _ := builder.FromFile("../tests_data/collection-01.json", BuilderOptions{})

	// Then
	if col == nil {
		t.Error("Cannot test extracting structures definitions, collection is nil")
		return
	}
	if reflect.DeepEqual(col.Structures, expectedStructures) == false {
		t.Errorf("Collection structures definition were not properly extracted, expected %v, got %v",
			expectedStructures, col.Structures)
	}
}
