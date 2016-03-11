package postman

type Folder struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Order       []string    `json:"order"`
	Owner       interface{} `json:"owner"`
}
