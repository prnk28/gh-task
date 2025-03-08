/*
Copyright © 2025 Prad N i@prad.nu
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/prnk28/gh-task/internal/ctx"
	"github.com/prnk28/gh-task/internal/ghc"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "GitHub CLI extension for task management",
	Long: `gh-task is a GitHub CLI extension that helps manage tasks across repositories.

It integrates with organization-level Taskfiles to provide standardized task execution
across repositories within an organization. The extension requires a .github repository
in the organization to store shared task configurations.

Usage:
  gh task [command]`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		c, err := ctx.Get(cmd)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		ok := ghc.OrgHasRepo(c.Current.RepoOwner, ".github")
		if !ok {
			cmd.Println("gh-task: .github repo required")
			return
		}
		taskfilePath, err := c.GetTaskfile()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		// Create the task command with the taskfile flag
		taskCmd := exec.Command("task", append([]string{"--taskfile", taskfilePath}, args...)...)
		
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
