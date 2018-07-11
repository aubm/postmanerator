package utils_test

import (
	"fmt"
	"os"
	"path"

	. "github.com/aubm/postmanerator/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Git", func() {

	var (
		gitAgent *GitAgent
	)

	BeforeEach(func() {
		gitAgent = &GitAgent{}
	})

	Describe("Clone", func() {

		var (
			cloneArgsToUse    []string
			cloneOptionsToUse CloneOptions
			returnedError     error
		)

		JustBeforeEach(func() {
			returnedError = gitAgent.Clone(cloneArgsToUse, cloneOptionsToUse)
		})

		BeforeEach(func() {
			workingDir := path.Join(os.TempDir(), uuid.NewV4().String())
			cloneOptionsToUse.TargetDirectory = workingDir
			must(os.Mkdir(cloneOptionsToUse.TargetDirectory, 0777))
		})

		AfterEach(func() {
			must(os.RemoveAll(cloneOptionsToUse.TargetDirectory))
		})

		Context("when everything is ok", func() {

			BeforeEach(func() {
				cloneArgsToUse = []string{"https://github.com/aubm/postmanerator-default-theme", "default"}
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should create a directory with the theme files", func() {
				Expect(fmt.Sprintf("%s/default", cloneOptionsToUse.TargetDirectory)).To(BeADirectory())
				Expect(fmt.Sprintf("%s/default/index.tpl", cloneOptionsToUse.TargetDirectory)).To(BeAnExistingFile())
			})

		})

		Context("when the url is not valid", func() {

			BeforeEach(func() {
				cloneArgsToUse = []string{"foo"}
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(ContainSubstring("Git clone failed:"))
			})

		})

	})

})
