package themes_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestThemes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Themes Suite")
}

func readFileContent(filePath string) string {
	f := mustReturn(os.Open(filePath)).(*os.File)
	defer f.Close()
	b := mustReturn(ioutil.ReadAll(f)).([]byte)
	return string(b[:])
}

func mustReturn(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}

func must(v ...interface{}) {
	err := v[len(v)-1]
	if err != nil {
		panic(err)
	}
}
