package ctx

import (
	"os"
	"path/filepath"

	"github.com/prnk28/gh-task/internal/ghc"
)

func getAppConfigHome() (string, error) {
	xdgHome, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(xdgHome, "gh-task"), nil
}

func getOrgTaskfilesHome(org string) (string, error) {
	confHome, err := getAppConfigHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(confHome, "src", org), nil
}

func mkDirOrg(org string) (string, error) {
	taskfilesDir, err := getOrgTaskfilesHome(org)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(taskfilesDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return taskfilesDir, nil
}

func downloadOrgData(org string) (string, error) {
	exists, path, err := orgDirExists(org)
	if err != nil {
		return "", err
	}
	if exists {
		return path, nil
	}

	if !exists {
		// 1. Create taskfiles directory for org
		dlDir, err := mkDirOrg(org)
		if err != nil {
			return "", err
		}

		// 2. Download Taskfile.yml
		out, err := ghc.QueryDownloadFile(org, "Taskfile.yml", dlDir).Exec()
		if err != nil {
			return "", err
		}

		// 3. Download taskfiles directory
		_, err = ghc.QueryDownloadFolder(org, "taskfiles", dlDir).Exec()
		if err != nil {
			return "", err
		}
		return out, nil
	}
	return path, nil
}

func orgDirExists(org string) (bool, string, error) {
	home, _ := getOrgTaskfilesHome(org)
	_, err := os.Stat(home)
	if err != nil {
		if os.IsNotExist(err) {
			return false, home, nil
		}
		return false, "", err
	}
	return true, home, nil
}

func rmDirOrg(org string) error {
	taskfilesDir, err := getOrgTaskfilesHome(org)
	if err != nil {
		return err
	}
	return os.RemoveAll(taskfilesDir)
}
