package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	subCmdName = "sub"

	subCmdFlagNameInput           = "input"
	subCmdFlagNameInputShort      = "i"
	subCmdFlagNameExpression      = "expression"
	subCmdFlagNameExpressionShort = "e"
	subCmdFlagNameOutput          = "output"
	subCmdFlagNameOutputShort     = "o"
)

var subCmd = &cobra.Command{
	Use:   subCmdName + " [flags] [...substitutes]",
	Short: "Regex replace in strings or files",
	Long: subCmdName + ` is a small utility that uses regex and regex groups to substitute
text inside either a piped in string or a file with arbitrarily numbered inputs.
The results are printed via stdout.

The regex flavor is Go, which contains group names. In the context of ` + subCmdName + `
this allows specifying exact replacement strings within the CLI args by placing
the strings at the correct CLI arg index. Here's a simple example:

	> goblin echo Hello World | goblin ` + subCmdName + ` -e "Hello (?P<1>(?s).*)" "Foo"
	Hello Foo

	> goblin echo Hello World | goblin ` + subCmdName + ` -e "(?P<2>Hello) (?P<1>World)" "Foo" "Bar"
	Bar Foo

` + subCmdName + ` groups are labelled sequentially with the [...substitutes] passed with
group 0 specifically referring to stdin/piped data:

	> goblin echo Hello World | goblin ` + subCmdName + ` -e "Hello(?P<0>\s)World"
	HelloHello WorldWorld

File data can also be loaded as a haystack by using the --input flag. If this flag is used then
stdin is not operated on as a haystack:

	README.md:
		Goblin is pretty cool bro.

	> goblin echo deez nutz | goblin ` + subCmdName + ` -i README.md -e "(?P<0>is pretty cool) -o README.md"

	README.md:
		Goblin deez nutz bro.
`,
	Run: func(cmd *cobra.Command, args []string) {
		stdinBuff, stdinReadErr := io.ReadAll(os.Stdin)
		stdinStr := string(stdinBuff)

		// Fetch the appropriate input
		operationStr := ""

		inputFileSet := cmd.Flags().Lookup(subCmdFlagNameInput).Changed
		if inputFileSet {
			inputFilePath, _ := cmd.Flags().GetString(subCmdFlagNameInput)
			fileData, err := os.ReadFile(inputFilePath)
			if err != nil {
				log.Fatal(err)
				return
			}
			operationStr = string(fileData)
		} else {
			if stdinReadErr != nil {
				log.Fatal(stdinReadErr)
			}
			operationStr = stdinStr
		}

		expression, _ := cmd.Flags().GetString("expression")

		// Compile regex
		regex, err := regexp.Compile(expression)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Search
		match := regex.FindStringSubmatchIndex(operationStr)
		if len(match) == 0 {
			log.Fatal("no matches for regex input")
		}

		totalMatches := (len(match) - 2) / 2
		expNames := regex.SubexpNames()[1:]

		// Check if the named groups matches the amount of matches we've made
		if totalMatches != len(expNames) {
			log.Fatal(fmt.Sprintf("the amount of named groups (%d) does not equal the amount of matches made (%d)", len(expNames), totalMatches))
		}

		// Strip the first 2 entries, since these are not group related
		match = match[2:]
		// Prepend 0 (start of operator string) to the beginning of the match
		match = append([]int{0}, match...)
		// Append length of total string to end of the match
		match = append(match, len(operationStr))

		// Prepend stdin to the groups if there is one piped in
		args = append([]string{stdinStr}, args...)

		var outputArr []string
		for i, name := range expNames {
			// Don't allow nameless groups
			if name == "" {
				log.Fatal(fmt.Sprintf("group %d is not named", i+1))
			}

			// Convert the name to an integer
			groupNum, err := strconv.Atoi(name)
			if err != nil {
				log.Fatal(fmt.Sprintf("group %d (%s) name cannot be converted to an integer", i+1, name))
			}

			// Check if there is input for this group number
			if groupNum > len(args)-1 {
				log.Fatal(fmt.Sprintf("cannot reference argument %d from group %d", groupNum, i+1))
			}

			strIdxStart := match[i*2]
			strIdxEnd := match[i*2+1]
			outputArr = append(outputArr, operationStr[strIdxStart:strIdxEnd])
			outputArr = append(outputArr, args[groupNum])
		}

		// Attach the end
		beginLastIdx := match[len(match)-2]
		outputArr = append(outputArr, operationStr[beginLastIdx:])

		// Combine the output
		output := strings.Join(outputArr[:], "")
		f := os.Stdout

		isOutputSet := cmd.Flags().Lookup(subCmdFlagNameOutput).Changed
		if isOutputSet {
			outputFilePath, err := cmd.Flags().GetString(subCmdFlagNameOutput)
			if outputFilePath != "" {
				f, err = os.Create(outputFilePath)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		f.WriteString(output)
	},
}

func init() {
	rootCmd.AddCommand(subCmd)

	subCmd.Flags().StringP(subCmdFlagNameExpression, subCmdFlagNameExpressionShort, "", "regex expression to use as a needle")
	subCmd.Flags().StringP(subCmdFlagNameInput, subCmdFlagNameInputShort, "", "use file contents as haystack. Not specifying this flag will cause stdin to be the\nhaystack and regex group <0> to be available as a replacement group")
	subCmd.Flags().StringP(subCmdFlagNameOutput, subCmdFlagNameOutputShort, "", "output file to write the results to")

	subCmd.MarkFlagRequired(subCmdFlagNameExpression)
}
