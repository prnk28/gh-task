package ctx

import (
	"strings"

	"github.com/prnk28/gh-task/internal/ghc"
)

var peerExtensions = []string{
	"yuler/gh-download",
}

func checkPeerDeps() map[string]bool {
	res := make(map[string]bool)
	out, err := ghc.CmdArgs("extension", "list").Exec()
	if err != nil {
		return res
	}

	for _, ext := range peerExtensions {
		if strings.Contains(out, ext) {
			res[ext] = true
			continue
		} else {
			res[ext] = false
			ghc.CmdArgs("extension", "install", ext).Exec()
		}
	}
	return res
}
