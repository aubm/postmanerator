package themes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestThemes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Themes Suite")
}

func must(v ...interface{}) {
	err := v[len(v)-1]
	if err != nil {
		panic(err)
	}
}
