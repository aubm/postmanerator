package helper

import (
	"reflect"
	"testing"

	"github.com/aubm/postmanerator/postman"
)

func TestFindRequest(t *testing.T) {
	requests := []postman.Request{
		{ID: "azerty", URL: "http://{{domain}}/api/chats"},
		{ID: "querty", URL: "http://{{domain}}/api/cats"},
	}

	cases := []struct {
		in  string
		req *postman.Request
	}{
		{"azerty", &requests[0]},
		{"foo", nil},
		{"querty", &requests[1]},
		{"bar", nil},
	}

	for i := 0; i < len(cases); i++ {
		req := findRequest(requests, cases[i].in)
		if !reflect.DeepEqual(req, cases[i].req) {
			t.Errorf("when req id = %v, expected req to equal %v, got %v", cases[i].in, cases[i].req, req)
		}
	}
}

func TestCurlSnippet(t *testing.T) {
	requests := []postman.Request{
		{RawHeaders: "Content-Type: application/json\nAccept: */*\n", URL: "http://{{domain}}/api/items",
			Method: "POST", RawModeData: "{\n    \"foo\": \"bar\"\n}", DataMode: "raw"},
		{RawHeaders: "Accept: */*\n", URL: "http://{{domain}}/api/items/45",
			Method: "DELETE", RawModeData: "", DataMode: "raw"},
		{RawHeaders: "Accept: */*\n", URL: "http://{{domain}}/api/items/45",
			Method: "GET", RawModeData: "some data", DataMode: "raw"},
		{RawHeaders: "Accept: */*\n", URL: "http://{{domain}}/api/items",
			Method: "POST", Data: []postman.RequestData{{"firstname", "foo"}, {"lastname", "bar"}}, DataMode: "urlencoded"},
		{RawHeaders: "Accept: */*\n", URL: "http://{{domain}}/api/items",
			Method: "POST", Data: []postman.RequestData{{"firstname", "foo"}, {"lastname", "bar"}}, DataMode: "params"},
	}
	expectedOutputs := []string{
		`curl -X POST -H "Content-Type: application/json" -H "Accept: */*" -d '{
    "foo": "bar"
}' "http://{{domain}}/api/items"`,
		`curl -X DELETE -H "Accept: */*" "http://{{domain}}/api/items/45"`,
		`curl -X GET -H "Accept: */*" "http://{{domain}}/api/items/45"`,
		`curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Accept: */*" -d "firstname=foo&lastname=bar" "http://{{domain}}/api/items"`,
		`curl -X POST -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW" -H "Accept: */*" -F "firstname=foo" -F "lastname=bar" "http://{{domain}}/api/items"`,
	}

	for i, req := range requests {
		output := curlSnippet(req)
		if output != expectedOutputs[i] {
			t.Errorf("when i = %v, expected curl snippet to be:\n%v\n\nbut got:\n%v", i, expectedOutputs[i], output)
		}
	}
}
