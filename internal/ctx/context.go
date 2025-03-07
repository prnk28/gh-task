package ctx

import (
	"context"
	"fmt"

	"github.com/cli/go-gh"
	"github.com/spf13/cobra"
)

// Define a key type for context values
type contextKey string

// Define the key for storing the Context in the cobra command's context
const ctxKey = contextKey("gh-pm-context")

type Context struct {
	Orgs     []string        `json:"orgs"`
	Name     string          `json:"name"`
	Login    string          `json:"login"`
	Current  Current         `json:"current"`
	PeerDeps map[string]bool `json:"peerdeps"`
}

func (c *Context) String() string {
	return fmt.Sprintf("Context{Orgs: %v, Name: %v, Login: %v, Current: %v, PeerDeps: %v}", c.Orgs, c.Name, c.Login, c.Current, c.PeerDeps)
}

type Current struct {
	RepoName  string `json:"repo_name"`
	RepoOwner string `json:"repo_owner"`
	Branch    string `json:"branch"`
	Path      string `json:"path"`
}

func (c *Current) String() string {
	return fmt.Sprintf("Current{RepoName: %v, RepoOwner: %v, Branch: %v, Path: %v}", c.RepoName, c.RepoOwner, c.Branch, c.Path)
}

func Get(cmd *cobra.Command) (*Context, error) {
	// Try to retrieve existing context
	cmdCtx := cmd.Context()
	if cmdCtx == nil {
		cmdCtx = context.Background()
		cmd.SetContext(cmdCtx)
	}

	// Check if context already exists
	if existingCtx, ok := cmdCtx.Value(ctxKey).(*Context); ok {
		return existingCtx, nil
	}

	// Create new context if it doesn't exist
	orgs, err := listOrgs()
	if err != nil {
		return nil, err
	}

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
	deps := checkPeerDeps()
	newCtx := &Context{
		Name:  repo.Name(),
		Login: repo.Owner(),
		Current: Current{
			RepoName:  repo.Name(),
			RepoOwner: repo.Owner(),
			Path:      wrkDir,
			Branch:    branch,
		},
		Orgs:     orgs,
		PeerDeps: deps,
	}

	// Create a new context with our value
	updatedCtx := context.WithValue(cmdCtx, ctxKey, newCtx)
	cmd.SetContext(updatedCtx)
	return newCtx, nil
}
