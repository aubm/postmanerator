package main

import (
	"errors"
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
		req postman.Request
		err error
	}{
		{"azerty", requests[0], nil},
		{"foo", postman.Request{}, errors.New("request not found")},
		{"querty", requests[1], nil},
		{"bar", postman.Request{}, errors.New("request not found")},
	}

	for i := 0; i < len(cases); i++ {
		req, err := findRequest(requests, cases[i].in)
		if !reflect.DeepEqual(req, cases[i].req) {
			t.Errorf("when req id = %v, expected req to equal %v, got %v", cases[i].in, cases[i].req, req)
		}
		if !reflect.DeepEqual(err, cases[i].err) {
			t.Errorf("when req id = %v, expected err to equal %v, got %v", cases[i].in, cases[i].err, err)
		}
	}
}
