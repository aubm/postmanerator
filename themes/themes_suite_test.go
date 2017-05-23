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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
