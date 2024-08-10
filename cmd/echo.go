package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	echoCmdName = "echo"

	echoCmdFlagName      = "tofile"
	echoCmdFlagNameShort = "f"
)

var echoCmd = &cobra.Command{
	Use:   echoCmdName + " [flags] [message]",
	Short: "Prints a message, without a newline and without quotes",
	Long: `Why do we need an echo command when every OS has one already? Cross-platformness is the answer,
along with some weirdness on Windows. Windows can goblin deez nutz.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := os.Stdout

		// If file is specified to dump to, then use that as the fd
		toFileSet := cmd.Flags().Lookup(echoCmdFlagName).Changed
		if toFileSet {
			filePath := cmd.Flags().Lookup(echoCmdFlagName).Value.String()
			var err error
			f, err = os.Create(filePath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Check stdin first, if there's something there then don't print args
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			buf, _ := io.ReadAll(os.Stdin)
			f.WriteString(string(buf))
		} else {
			// Dump all input arguments directly to stdout, without a newline, of course
			f.WriteString(strings.Join(args[:], " "))
		}

		if toFileSet {
			f.Close()
		}
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)

	echoCmd.Flags().StringP(echoCmdFlagName, echoCmdFlagNameShort, "", "regex expression to use as a needle")
}
