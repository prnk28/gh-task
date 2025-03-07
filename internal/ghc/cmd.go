package ghc

import (
	"encoding/json"
	"strings"

	"github.com/cli/go-gh"
)

func Cmd(s string) GHCommand {
	clean := strings.ReplaceAll(s, "gh ", "")
	ptrs := strings.Split(clean, " ")
	return GHCommand(ptrs)
}

func CmdArgs(args ...string) GHCommand {
	return GHCommand(args)
}

type GHCommand []string

func (c GHCommand) Exec() (string, error) {
	out, _, err := gh.Exec(c.StringArray()...)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// ExecUnmarshal unmarshals the output of the command into the provided interface with JSON
func (c GHCommand) ExecUnmarshal(i any) error {
	out, err := c.Exec()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(out), i)
}

// StringArray returns the command as an array of strings
func (c GHCommand) StringArray() []string {
	return []string(c)
}
