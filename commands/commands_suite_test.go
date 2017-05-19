package commands_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

func putFileContents(filePath, contents string) {
	f := createFile(filePath)
	defer f.Close()
	must(f.Write([]byte(contents)))
}

func readFileContents(filePath string) string {
	f := openFile(filePath)
	defer f.Close()
	b := must(ioutil.ReadAll(f)).([]byte)
	return string(b[:])
}

func createFile(filePath string) *os.File {
	return must(os.Create(filePath)).(*os.File)
}

func openFile(filePath string) *os.File {
	return must(os.Open(filePath)).(*os.File)
}

func must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
