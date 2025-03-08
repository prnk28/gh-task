/*
Copyright © 2025 Prad N i@prad.nu
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prnk28/gh-task/internal/ctx"
	"github.com/prnk28/gh-task/internal/ghc"
	"github.com/spf13/cobra"
)

var (
	cfgFile            string
	printPath          bool
	cachedTaskfilePath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Gh CLI extension for global organization-level taskfile execution",
	Long: `gh-task is a GitHub CLI extension that helps consolidate taskfiles across repositories.
	It integrates with organization-level Taskfiles to provide standardized task execution
	across repositories within an organization. The extension requires a .github repository
	in the organization to store shared task configurations.

	Usage:
  	gh task [command]`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Get the taskfile path
		taskfilePath, err := getTaskfilePath(cmd)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		// If print-path flag is set, just print the path and return
		if printPath {
			cmd.Println(taskfilePath)
			return
		}

		// Create the task command with the taskfile flag
		currentPath, err := os.Getwd()
		if err != nil {
			cmd.PrintErrf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Create the task command with the taskfile flag
		taskCmd := exec.Command("task", append([]string{"--taskfile", taskfilePath, "--dir", currentPath}, args...)...)

		// Set up pipes for stdin, stdout, and stderr
		taskCmd.Stdin = os.Stdin
		taskCmd.Stdout = os.Stdout
		taskCmd.Stderr = os.Stderr

		// Run the task command
		err = taskCmd.Run()
		if err != nil {
			cmd.PrintErrf("Error executing task: %v\n", err)
			// If the command fails with a specific exit code, propagate it
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&printPath, "print-path", "p", false, "Print the path to the Taskfile instead of executing it")
}

// getTaskfilePath gets the taskfile path with caching for better performance
func getTaskfilePath(cmd *cobra.Command) (string, error) {
	// Return cached path if available
	if cachedTaskfilePath != "" {
		// Verify the file still exists
		if _, err := os.Stat(cachedTaskfilePath); err == nil {
			return cachedTaskfilePath, nil
		}
	}

	c, err := ctx.Get(cmd)
	if err != nil {
		return "", err
	}

	// Quick check if the organization has a .github repo
	ok := ghc.OrgHasRepo(c.Current.RepoOwner, ".github")
	if !ok {
		return "", fmt.Errorf("gh-task: .github repo required for organization %s", c.Current.RepoOwner)
	}

	// Get the expected taskfile path
	taskfilePath, err := c.GetTaskfile()
	if err != nil {
		return "", err
	}

	// Check if the taskfile directory exists before trying to access it
	taskfileDir := filepath.Dir(taskfilePath)
	if _, err := os.Stat(taskfileDir); os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(taskfileDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create taskfile directory: %w", err)
		}
	}

	// Check if the taskfile exists
	if _, err := os.Stat(taskfilePath); os.IsNotExist(err) {
		// Taskfile doesn't exist, try to download it
		orgDir, err := ctx.DownloadOrgData(c.Current.RepoOwner)
		if err != nil {
			return "", fmt.Errorf("failed to download organization data: %w", err)
		}

		// Check again after download
		if _, err := os.Stat(taskfilePath); os.IsNotExist(err) {
			return "", fmt.Errorf("taskfile not found at %s even after download", taskfilePath)
		}
		cachedTaskfilePath = orgDir
	}

	// Cache the path for future use
	cachedTaskfilePath = taskfilePath
	return taskfilePath, nil
}

// func init() {
// 	cobra.OnInitialize(initConfig)
//
// 	// Here you will define your flags and configuration settings.
// 	// Cobra supports persistent flags, which, if defined here,
// 	// will be global for your application.
//
// 	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-task.yaml)")
//
// 	// Cobra also supports local flags, which will only run
// 	// when this action is called directly.
// 	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
//
// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := os.UserHomeDir()
// 		cobra.CheckErr(err)
//
// 		// Search config in home directory with name ".gh-task" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigType("yaml")
// 		viper.SetConfigName(".gh-task")
// 	}
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
// 	}
// }
