package postman

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestEnvironmentFromFile(t *testing.T) {
	filename := createTmpEnvironmentFile()

	env, err := EnvironmentFromFile(filename)

	expectedEnv := map[string]string{"domain": "localhost:8080"}
	if ok := reflect.DeepEqual(env, expectedEnv); ok == false {
		t.Errorf("Expected %v, got %v, err is %v", expectedEnv, env, err)
	}
}

func createTmpEnvironmentFile() string {
	f, err := ioutil.TempFile(os.TempDir(), "postman_env_file")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprint(f, `{
	"id": "316cfffe-80bc-ff35-4a30-46d6085d1973",
	"name": "Books API - Local",
	"values": [
		{
			"key": "domain",
			"value": "localhost:8080",
			"type": "text",
			"enabled": true
		}
	],
	"timestamp": 1476905519649,
	"_postman_variable_scope": "environment",
	"_postman_exported_at": "2016-10-19T19:33:45.473Z",
	"_postman_exported_using": "Postman/4.7.2"
}`)

	return f.Name()
}
