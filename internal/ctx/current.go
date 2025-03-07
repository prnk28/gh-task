package ctx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cli/go-gh"
)

var ghcliExtensionDeps = []string{
	"valeriobelli/gh-milestone",
}

type Current struct {
	RepoName  string          `json:"repo_name"`
	RepoOwner string          `json:"repo_owner"`
	Branch    string          `json:"branch"`
	Path      string          `json:"path"`
	PeerDeps  map[string]bool `json:"peerdeps"`
	Taskfile  string          `json:"taskfile"`
}

func (c *Current) String() string {
	return fmt.Sprintf("Current{RepoName: %v, RepoOwner: %v, Branch: %v, Path: %v, PeerDeps: %v, Taskfile: %v}", c.RepoName, c.RepoOwner, c.Branch, c.Path, c.PeerDeps, c.Taskfile)
}

func fetchCurrent() (*Current, error) {
	repo, err := gh.CurrentRepository()
	if err != nil {
		return nil, err
	}
	wrkDir, err := WorkingDir()
	if err != nil {
		return nil, err
	}
	wrkBranch, err := CurrentBranch()
	branch := ""
	if err == nil {
		branch = wrkBranch
	}
	currRepo := repo.Name()
	currOrg := repo.Owner()
	p, err := downloadOrgTaskfileData(currOrg)
	if err != nil {
		return nil, err
	}
	return &Current{
		RepoName:  currRepo,
		RepoOwner: currOrg,
		Path:      wrkDir,
		Branch:    branch,
		PeerDeps:  checkPeerDeps(),
		Taskfile:  p,
	}, nil
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
