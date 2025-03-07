package ctx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Define a key type for context values
type contextKey string

// Define the key for storing the Context in the cobra command's context
const ctxKey = contextKey("gh-pm-context")

type Context struct {
	ConfigHome string   `json:"config_home"`
	Orgs       []string `json:"orgs"`
	Name       string   `json:"name"`
	Login      string   `json:"login"`
	Current    *Current `json:"current"`
}

func initContext(cmdCtx context.Context, cmd *cobra.Command, curr *Current) (*Context, error) {
	// Fetch XDG config home
	configHome, err := getAppConfigHome()
	if err != nil {
		return nil, err
	}

	// Create config home if it doesn't exist
	err = os.MkdirAll(configHome, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create new context if it doesn't exist
	orgs, err := filterActiveOrgs()
	if err != nil {
		return nil, err
	}

	newCtx := &Context{
		ConfigHome: configHome,
		Name:       curr.RepoName,
		Login:      curr.RepoOwner,
		Orgs:       orgs,
		Current:    curr,
	}

	// Create a new context with our value
	updatedCtx := context.WithValue(cmdCtx, ctxKey, newCtx)
	cmd.SetContext(updatedCtx)
	return newCtx, err
}

func (c *Context) String() string {
	return fmt.Sprintf("Context{Orgs: %v, Name: %v, Login: %v, Current: %v}", c.Orgs, c.Name, c.Login, c.Current)
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
		fmt.Println("Using existing context")
		return existingCtx, nil
	}

	curr, err := fetchCurrent()
	if err != nil {
		return nil, err
	}

	newCtx, err := initContext(cmdCtx, cmd, curr)
	if err != nil {
		return nil, err
	}

	// Create a new context with our value
	updatedCtx := context.WithValue(cmdCtx, ctxKey, newCtx)
	cmd.SetContext(updatedCtx)
	return newCtx, nil
}

func (c *Context) GetTaskfile() (string, error) {
	return filepath.Join(c.ConfigHome, "src", c.Current.RepoOwner, "Taskfile.yml"), nil
}
