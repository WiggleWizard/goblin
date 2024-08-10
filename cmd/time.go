package cmd

import (
	"github.com/itchyny/timefmt-go"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Prints the current date and time to stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		formatSet := cmd.Flags().Lookup("format").Changed
		unixSet := cmd.Flags().Lookup("unix").Changed

		if formatSet {
			// Format output
			format, err := cmd.Flags().GetString("format")
			if err == nil {
				t := timefmt.Format(time.Now(), format)
				log.Print(t)
				return
			}
		} else if unixSet {
			// Unix timestamp
			log.Print(time.Now().Unix())
			return
		}

		// Spit out help if no switches were valid
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().BoolP("unix", "u", false, "outputs the date and time in Unix timestamp")
	timeCmd.Flags().StringP("format", "f", "", "see https://man7.org/linux/man-pages/man3/strftime.3.html for formatting")
}
