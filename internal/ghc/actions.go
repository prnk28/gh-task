package ghc

import (
	"slices"
	"strings"
)

func OrgHasRepo(org, repo string) bool {
	out, err := QueryOrgRepos(org).Exec()
	if err != nil {
		return false
	}
	repos := strings.Split(out, "\n")
	return slices.Contains(repos, repo)
}
