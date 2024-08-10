package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var genusageCmd = &cobra.Command{
	Use:   "genusage",
	Short: "Generates a total list of all commands, usages and descriptions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		genAsMarkdown := cmd.Flags().Lookup("md").Changed

		if genAsMarkdown {
			for _, subcommand := range rootCmd.Commands() {
				if subcommand.Hidden {
					continue
				}
				if subcommand.Name() == "help" {
					continue
				}

				log.Println("## " + subcommand.Name())
				log.Println(subcommand.Short)
				log.Println("```")
				if len(subcommand.Long) > 0 {
					log.Println(subcommand.Long + "\n")
				}
				log.Print(subcommand.UsageString())
				log.Println("```")
			}
		} else {
			for _, subcommand := range rootCmd.Commands() {
				log.Println(subcommand.UsageString())
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(genusageCmd)

	genusageCmd.Flags().Bool("md", false, "generate each command usage inside a markdown code block")
}
