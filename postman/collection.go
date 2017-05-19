package postman

type Collection struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Order       []string    `json:"order"`
	Folders     []Folder    `json:"folders"`
	Timestamp   int64       `json:"timestamp"`
	Owner       interface{} `json:"owner"`
	RemoteLink  string      `json:"remoteLink"`
	Public      bool        `json:"public"`
	Requests    []Request   `json:"requests"`
	Structures  []StructureDefinition
}
