package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

type GitAgent struct{}

func (c GitAgent) Clone(args []string, options CloneOptions) error {
	args = append([]string{"clone"}, args...)
	cmd := exec.Command("git", args...)
	cmd.Dir = options.TargetDirectory
	stdErr := new(bytes.Buffer)
	cmd.Stderr = stdErr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Git clone failed: %v", stdErr.String())
	}
	return nil
}

type CloneOptions struct {
	TargetDirectory string
}
