package postman

type Collection struct {
	Name        string
	Description string
	Requests    []Request
	Folders     []Folder
	Structures  []StructureDefinition
}

type Request struct {
	ID            string
	Name          string
	Description   string
	Method        string
	URL           string
	PayloadType   string
	PayloadRaw    string
	PayloadParams []KeyValuePair
	PathVariables []KeyValuePair
	Headers       []KeyValuePair
	Responses     []Response
}

type Response struct {
	ID         string
	Name       string
	Status     string
	StatusCode int
	Body       string
	Headers    []KeyValuePair
}

type Folder struct {
	ID          string
	Name        string
	Description string
	Requests    []Request
}

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

type KeyValuePair struct {
	Name        string
	Key         string
	Value       interface{}
	Description string
}
