package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type cmdDescriptor struct {
	Name string
}

var rootCmd = &cobra.Command{
	Use:   "goblin",
	Short: "A small utility to assist with Makefiles, CI or just scripting",
	Long:  ``,
	// Disable completion command
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
