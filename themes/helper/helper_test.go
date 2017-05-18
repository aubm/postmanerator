package helper

import (
	"reflect"
	"testing"

	"github.com/aubm/postmanerator/postman"
)

func TestFindRequest(t *testing.T) {
	requests := []postman.Request{
		{ID: "azerty", URL: "http://foo.bar/api/chats"},
		{ID: "querty", URL: "http://foo.bar/api/cats"},
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

var stubRequests []postman.Request = []postman.Request{
	{RawHeaders: "Content-Type: application/json\nAccept: */*\n", URL: "http://foo.bar/api/items",
		Method: "POST", RawModeData: "{\n    \"foo\": \"bar\"\n}", DataMode: "raw"},
	{RawHeaders: "Accept: */*\n", URL: "http://foo.bar/api/items/45",
		Method: "DELETE", RawModeData: "", DataMode: "raw"},
	{RawHeaders: "Accept: */*\n", URL: "http://foo.bar/api/items/45",
		Method: "GET", RawModeData: "some data", DataMode: "raw"},
	{RawHeaders: "Accept: */*\n", URL: "http://foo.bar/api/items",
		Method: "POST", Data: []postman.RequestData{{"firstname", "foo"}, {"lastname", "bar"}}, DataMode: "urlencoded"},
	{RawHeaders: "Accept: */*\n", URL: "http://foo.bar/api/items",
		Method: "POST", Data: []postman.RequestData{{"firstname", "foo"}, {"lastname", "bar"}}, DataMode: "params"},
}

func TestCurlSnippet(t *testing.T) {
	expectedOutputs := []string{
		`curl -X POST -H "Content-Type: application/json" -H "Accept: */*" -d '{
    "foo": "bar"
}' "http://foo.bar/api/items"`,
		`curl -X DELETE -H "Accept: */*" "http://foo.bar/api/items/45"`,
		`curl -X GET -H "Accept: */*" "http://foo.bar/api/items/45"`,
		`curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Accept: */*" -d "firstname=foo&lastname=bar" "http://foo.bar/api/items"`,
		`curl -X POST -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW" -H "Accept: */*" -F "firstname=foo" -F "lastname=bar" "http://foo.bar/api/items"`,
	}

	for i, req := range stubRequests {
		output := curlSnippet(req)
		if output != expectedOutputs[i] {
			t.Errorf("when i = %v, expected curl snippet to be:\n%v\n\nbut got:\n%v", i, expectedOutputs[i], output)
		}
	}
}

func TestHttpSnippet(t *testing.T) {
	expectedOutputs := []string{
		`POST /api/items HTTP/1.1
Host: foo.bar
Content-Type: application/json
Accept: */*

{
    "foo": "bar"
}`,
		`DELETE /api/items/45 HTTP/1.1
Host: foo.bar
Accept: */*`,
		`GET /api/items/45 HTTP/1.1
Host: foo.bar
Accept: */*`,
		`POST /api/items HTTP/1.1
Host: foo.bar
Accept: */*
Content-Type: application/x-www-form-urlencoded

firstname=foo&lastname=bar`,
		`POST /api/items HTTP/1.1
Host: foo.bar
Accept: */*
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="firstname"

foo
----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="lastname"

bar
----WebKitFormBoundary7MA4YWxkTrZu0gW`,
	}

	for i, req := range stubRequests {
		output := httpSnippet(req)
		if output != expectedOutputs[i] {
			t.Errorf("when i = %v, expected http snippet to be:\n%v\n\nbut got:\n%v", i, expectedOutputs[i], output)
		}
	}
}

func TestSlugify(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello, World!", "hello-world"},
		{"Lorem @ Ipsum? Etc.", "lorem-ipsum-etc"},
		{"Vivement l'été", "vivement-l-t"},
	}

	for _, c := range cases {
		slug := slugify(c.in)
		if slug != c.out {
			t.Errorf("for in = %v, expected out to equal %v, but got %v", c.in, c.out, slug)
		}
	}
}
