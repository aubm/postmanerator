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
