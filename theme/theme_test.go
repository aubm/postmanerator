package theme

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockCloner struct {
	returnError bool
	lastArgs    []string
}

func (c *MockCloner) Clone(args []string) error {
	c.lastArgs = args
	if c.returnError {
		return errors.New("git clone failed")
	}
	return nil
}

func TestGitClone(t *testing.T) {
	git := &MockCloner{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `[submodule "default"]
	path = default
	url = https://github.com/aubm/postmanerator-default-theme.git
[submodule "foo"]
	path = foo
	url = https://github.com/aubm/postmanerator-foo-theme.git`)
	}))
	defer ts.Close()

	data := []struct {
		themeToDownload string
		destinationDir  string
		args            []string
		err             error
		cloneErr        bool
	}{
		{"foo", "", []string{"https://github.com/aubm/postmanerator-foo-theme.git", "foo"}, nil, false},
		{"foo", "my-foo", []string{"https://github.com/aubm/postmanerator-foo-theme.git", "my-foo"}, nil, false},
		{"bar", "", nil, errors.New("An error occured while trying to resolve theme 'bar': theme not found"), false},
		{"foo", "", []string{"https://github.com/aubm/postmanerator-foo-theme.git", "foo"},
			errors.New("There was an error while executing git clone: git clone failed"), true},
		{"http://myrepo.git", "", []string{"http://myrepo.git"}, nil, false},
		{"https://myrepo.git", "", []string{"https://myrepo.git"}, nil, false},
		{"git@myrepo.git", "", []string{"git@myrepo.git"}, nil, false},
	}
	for i, d := range data {
		git.returnError = d.cloneErr
		err := GitClone(d.themeToDownload, d.destinationDir, ts.URL, git)
		if reflect.DeepEqual(git.lastArgs, d.args) == false {
			t.Errorf("for i = %v, git cloner received wrong args: %v", i, git.lastArgs)
		}
		if reflect.DeepEqual(d.err, err) == false {
			t.Errorf("for i = %v, err should be %v, got %v", i, d.err, err)
		}
		git.lastArgs = nil
	}
}
