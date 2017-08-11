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

	builderOptions := BuilderOptions{
		EnvironmentVariables: Environment{
			"domain": "localhost",
			"catId":  "1",
			"dogId":  "1",
		},
	}

	testFiles := []string{
		"tests_data/collection-01.json",
		"tests_data/collection-02.json",
	}

	for _, testFile := range testFiles {
		// When
		col, err := builder.FromFile(testFile, builderOptions)

		// Then
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if reflect.DeepEqual(col.Structures, expectedStructures) == false {
			t.Errorf("Collection structures definition were not properly extracted, expected %v, got %v",
				expectedStructures, col.Structures)
		}
	}
}
