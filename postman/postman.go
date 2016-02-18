package postman

type Collection struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Order       []string  `json:"order"`
	Folders     []Folder  `json:"folders"`
	Timestamp   int       `json:"timestamp"`
	Owner       int       `json:"owner"`
	RemoteLink  string    `json:"remoteLink"`
	Public      bool      `json:"public"`
	Requests    []Request `json:"requests"`
}

type Folder struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Order       []string `json:"order"`
	Owner       int      `json:"owner"`
}

type Request struct {
	ID               string        `json:"id"`
	Headers          string        `json:"headers"`
	URL              string        `json:"url"`
	PreRequestScript string        `json:"preRequestScript"`
	PathVariables    interface{}   `json:"pathVariables"`
	Method           string        `json:"method"`
	Data             []interface{} `json:"data"`
	DataMode         string        `json:"dataMode"`
	Version          int           `json:"version"`
	Tests            string        `json:"tests"`
	CurrentHelper    string        `json:"currentHelper"`
	HelperAttributes interface{}   `json:"helperAttributes"`
	Time             int           `json:"time"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	CollectionID     string        `json:"collectionId"`
	Responses        []Response    `json:"responses"`
	RawModeData      string        `json:"rawModeData"`
}

type Response struct {
	Status       string `json:"status"`
	ResponseCode struct {
		Code   int    `json:"code"`
		Name   string `json:"name"`
		Detail string `json:"detail"`
	} `json:"responseCode"`
	Time    int `json:"time"`
	Headers []struct {
		Name        string `json:"name"`
		Key         string `json:"key"`
		Value       string `value`
		Description string `json:"description"`
	} `json:"headers"`
	Cookies     []interface{} `json:"cookies"`
	Mime        string        `json:"mime"`
	Text        string        `json:"text"`
	Language    string        `json:"language"`
	RawDataType string        `json:"text"`
	State       struct {
		Size string `json:"size"`
	} `json:"state"`
	PreviewType            string `json:"previewType"`
	SearchResultScrolledTo int    `json:"searchResultScrolledTo"`
	ForceNoPretty          bool   `json:"forceNoPretty"`
	Write                  bool   `json:"write"`
	Empty                  bool   `json:"empty"`
	Failed                 bool   `json:"failed"`
	IsSample               bool   `json:"isSample"`
	ScrollToResult         bool   `json:"scrollToResult"`
	RunTests               bool   `json:"runTests"`
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Request                struct {
		URL     string `json:"url"`
		Headers []struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
		} `json:"headers"`
		Data     string `json:"data"`
		Method   string `json:"method"`
		DataMode string `json:"dataMode"`
	} `json:"request"`
}
