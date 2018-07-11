package themes_test

import "github.com/aubm/postmanerator/postman"

var exampleCollection = postman.Collection{
	Name:        "My Collection",
	Description: "My awesome collection description.\n\nIt is very cool because:\n\n- foo\n- bar\n- fizz\n- buzz",
	Requests: []postman.Request{
		{
			ID:            "51f1ea36-418a-11e7-b895-6f2275a2594f",
			Name:          "Get all cats",
			Description:   "Cats are good for *health*!",
			Method:        "GET",
			URL:           "https://my-api/cats",
			PayloadType:   "params",
			PayloadRaw:    "",
			PayloadParams: nil,
			PathVariables: nil,
			Headers: []postman.KeyValuePair{
				{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
			},
			Responses: []postman.Response{
				{
					ID:         "593e6c1a-418a-11e7-89de-3fc82fddb918",
					Name:       "OK",
					Status:     "OK",
					StatusCode: 200,
					Body:       `[{"name":"Barney"}]`,
					Headers: []postman.KeyValuePair{
						{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
					},
				},
			},
		},
	},
	Folders: []postman.Folder{
		{
			ID:          "5edfce7a-418a-11e7-b415-2385135fef4a",
			Name:        "Everything about dogs",
			Description: "Dogs are so _funny_!",
			Requests: []postman.Request{
				{
					ID:            "64dc5c94-418a-11e7-9fd3-4f1a853549cb",
					Name:          "Get one dog by id",
					Description:   "Because that particular dog is very *special*!",
					Method:        "GET",
					URL:           "https://my-api/dogs/:id",
					PayloadType:   "params",
					PayloadRaw:    "",
					PayloadParams: nil,
					PathVariables: []postman.KeyValuePair{{Name: "id", Key: "id"}},
					Headers: []postman.KeyValuePair{
						{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
					},
					Responses: []postman.Response{
						{
							ID:         "6a1bbd76-418a-11e7-b95d-03ef405e2322",
							Name:       "OK",
							Status:     "OK",
							StatusCode: 200,
							Body:       `{"name":"Sam"}`,
							Headers: []postman.KeyValuePair{
								{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
							},
						},
					},
				},
				{
					ID:            "6e345bac-418a-11e7-b879-e7d2b4bf7a22",
					Name:          "Create a new dog",
					Description:   "The family is growing!",
					Method:        "POST",
					URL:           "https://my-api/dogs",
					PayloadType:   "raw",
					PayloadRaw:    `{"name":"Sam JR"}`,
					PayloadParams: nil,
					PathVariables: nil,
					Headers: []postman.KeyValuePair{
						{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
					},
					Responses: []postman.Response{
						{
							ID:         "7369c486-418a-11e7-8c9f-8fecb3df53c9",
							Name:       "OK",
							Status:     "Created",
							StatusCode: 201,
							Body:       `{"name":"Sam JR"}`,
							Headers: []postman.KeyValuePair{
								{Name: "New-Dog-ID", Key: "New-Dog-ID", Value: "sam-jr"},
								{Name: "Content-Type", Key: "Content-Type", Value: "application/json"},
							},
						},
					},
				},
				{
					ID:          "1f26facc-4560-11e7-86d4-b3f9d53f63d7",
					Name:        "Create a new dog with urlencoded values",
					Description: "The family is growing!",
					Method:      "POST",
					URL:         "https://my-api/dogs",
					PayloadType: "urlencoded",
					PayloadRaw:  "",
					PayloadParams: []postman.KeyValuePair{
						{Name: "name", Key: "name", Value: "Sam JR"},
					},
					PathVariables: nil,
				},
				{
					ID:          "1f26facc-4560-11e7-86d4-b3f9d53f63d7",
					Name:        "Create a new dog with form values",
					Description: "The family is growing!",
					Method:      "POST",
					URL:         "https://my-api/dogs",
					PayloadType: "params",
					PayloadRaw:  "",
					PayloadParams: []postman.KeyValuePair{
						{Name: "name", Key: "name", Value: "Sam JR"},
					},
					PathVariables: nil,
				},
			},
		},
	},
	Structures: []postman.StructureDefinition{
		{
			Name:        "Cat",
			Description: "Cat are felines",
			Fields: []postman.StructureFieldDefinition{
				{Name: "Name", Description: "The name of the cat", Type: "string"},
			},
		},
		{
			Name:        "Dog",
			Description: "Dog are canines",
			Fields: []postman.StructureFieldDefinition{
				{Name: "Name", Description: "The name of the dog", Type: "string"},
			},
		},
	},
}
