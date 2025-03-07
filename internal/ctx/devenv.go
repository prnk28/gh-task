package ctx

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var ghcliExtensionDeps = []string{
	"valeriobelli/gh-milestone",
}

func WorkingDir() (string, error) {
	return os.Getwd()
}

// CurrentBranch returns the name of the current git branch in the working directory
func CurrentBranch() (string, error) {
	// Create command to run "git branch"
	cmd := exec.Command("git", "branch", "--show-current")

	// Create buffer to capture output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Get output and trim whitespace
	branchName := strings.TrimSpace(out.String())
	return branchName, nil
}
