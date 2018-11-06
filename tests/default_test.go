package tests

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestDefaultCommand(t *testing.T) {
	testCases, err := prepareTestCases()
	if err != nil {
		t.Fatalf("failed to prepare test cases: %v", err)
	}

	for _, tc := range testCases {
		for _, expectedOutput := range tc.ExpectedOutputs {
			t.Run(fmt.Sprintf("case %s with theme %s", tc.Name, expectedOutput.Theme), func(t *testing.T) {
				tmpOutputFile, err := ioutil.TempFile("", "postmanerator_test_generated_output_*.txt")
				if err != nil {
					t.Fatalf("failed to create tmp output file: %v", err)
				}
				if err := tmpOutputFile.Close(); err != nil {
					t.Fatalf("failed to close generated tmp output file: %v", err)
				}
				cmd := exec.Command("postmanerator",
					"-collection",
					tc.CollectionPath,
					"-theme",
					path.Join("themes", expectedOutput.Theme),
					"-output",
					tmpOutputFile.Name(),
				)
				stdout := new(bytes.Buffer)
				stderr := new(bytes.Buffer)
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				if err := cmd.Run(); err != nil {
					t.Fatalf("failed to run postmanerator command: %v, stdout: %s, stderr: %v", err, stdout, stderr)
				}

				generatedOutput, err := ioutil.ReadFile(tmpOutputFile.Name())
				if err != nil {
					t.Fatalf("failed to read generated output: %v", err)
				}

				expectedOutput, err := ioutil.ReadFile(expectedOutput.ExpectedOutputFile)
				if err != nil {
					t.Fatalf("failed to read expected output file: %v", err)
				}

				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(string(generatedOutput), string(expectedOutput), false)

				for _, diff := range diffs {
					if diff.Type != diffmatchpatch.DiffEqual {
						t.Errorf("found new differences between generated and expected output files contents: %v", diff)
					}
				}
			})
		}
	}
}

func prepareTestCases() ([]TestCaseData, error) {
	testCases := make([]TestCaseData, 0)
	expectedOutputRegexp := regexp.MustCompile(`(.+)\.txt$`)

	cases, err := ioutil.ReadDir("cases")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cases dir")
	}
	for _, c := range cases {
		if !c.IsDir() {
			continue
		}

		tc := TestCaseData{
			Name:            c.Name(),
			CollectionPath:  path.Join("cases", c.Name(), "collection.json"),
			ExpectedOutputs: make([]TestCaseDataExpectedOutputData, 0),
		}

		expectedOutputsPath := path.Join("cases", tc.Name, "expected_outputs")
		expectedOutputs, err := ioutil.ReadDir(expectedOutputsPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read expected outputs for case %s", tc.Name)
		}

		for _, eo := range expectedOutputs {
			if eo.IsDir() {
				continue
			}

			n := eo.Name()
			if !expectedOutputRegexp.MatchString(n) {
				continue
			}

			matches := expectedOutputRegexp.FindStringSubmatch(n)
			tc.ExpectedOutputs = append(tc.ExpectedOutputs, TestCaseDataExpectedOutputData{
				Theme:              matches[1],
				ExpectedOutputFile: path.Join(expectedOutputsPath, n),
			})
		}

		testCases = append(testCases, tc)
	}

	return testCases, nil
}

type TestCaseData struct {
	Name            string
	CollectionPath  string
	ExpectedOutputs []TestCaseDataExpectedOutputData
}

type TestCaseDataExpectedOutputData struct {
	Theme              string
	ExpectedOutputFile string
}
