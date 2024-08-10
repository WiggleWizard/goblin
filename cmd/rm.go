package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

type rmCommandInfo struct {
	CmdInfo                cmdDescriptor
	FlagNameRecursive      string
	FlagNameRecursiveShort string
	FlagNameForce          string
	FlagNameForceShort     string
}

var rmCmdInfo = rmCommandInfo{
	CmdInfo:                cmdDescriptor{Name: "rm"},
	FlagNameRecursive:      "recursive",
	FlagNameRecursiveShort: "r",
	FlagNameForce:          "force",
	FlagNameForceShort:     "f",
}

var rmCmd = &cobra.Command{
	Use:   rmCmdInfo.CmdInfo.Name,
	Short: "remove files or directories",
	Long:  rmCmdInfo.CmdInfo.Name + ` removes each specified file. By default, it does not remove directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		recursiveSet, _ := cmd.Flags().GetBool(rmCmdInfo.FlagNameRecursive)
		forceSet, _ := cmd.Flags().GetBool(rmCmdInfo.FlagNameForce)

		// Iterate through all arbitrary args
		for _, filePath := range args {
			// Delete folders recursively if specified
			if recursiveSet {
				err := os.RemoveAll(filePath)
				if !forceSet && err != nil {
					log.Fatal(err)
				}
			} else {
				fStat, err := os.Stat(filePath)
				if !forceSet && err != nil {
					log.Fatal(err)
				}
				if !fStat.IsDir() {
					err := os.Remove(filePath)
					if !forceSet && err != nil {
						log.Fatal(err)
					}
				} else {
					if !forceSet {
						log.Fatalf("%s is a directory", filePath)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().BoolP(rmCmdInfo.FlagNameRecursive, rmCmdInfo.FlagNameRecursiveShort, false, "remove directories and their contents recursively")
	rmCmd.Flags().BoolP(rmCmdInfo.FlagNameForce, rmCmdInfo.FlagNameForceShort, false, "ignore nonexistent files and arguments, never prompt")
}
