package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type restocsourceCmdDescriptor struct {
	CmdInfo                 cmdDescriptor
	FlagNameResDir          string
	FlagNameResDirShort     string
	FlagNameDest            string
	FlagNameDestShort       string
	FlagNameNamespace       string
	FlagNameNamespaceShort  string
	FlagNameStringType      string
	FlagNameStringTypeShort string
}

var restocsourceCmdInfo = restocsourceCmdDescriptor{
	CmdInfo:                 cmdDescriptor{Name: "restocsource"},
	FlagNameResDir:          "sourcedir",
	FlagNameResDirShort:     "s",
	FlagNameDest:            "destination",
	FlagNameDestShort:       "d",
	FlagNameNamespace:       "namespace",
	FlagNameNamespaceShort:  "n",
	FlagNameStringType:      "stringtype",
	FlagNameStringTypeShort: "t",
}

var restocsourceCmd = &cobra.Command{
	Use:   restocsourceCmdInfo.CmdInfo.Name,
	Short: "Generates a C/C++ source and header that contains embedded resource data",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceDir, _ := cmd.Flags().GetString(restocsourceCmdInfo.FlagNameResDir)
		destSourcePath, _ := cmd.Flags().GetString(restocsourceCmdInfo.FlagNameDest)
		namespace, _ := cmd.Flags().GetString(restocsourceCmdInfo.FlagNameNamespace)
		stringType, _ := cmd.Flags().GetString(restocsourceCmdInfo.FlagNameStringType)

		// Check first if source folder even exist
		if !exists(resourceDir) {
			log.Fatal(fmt.Sprintf("resource directory %s does not exist", resourceDir))
		}

		// Try strip off any extensions from the destination source file
		sourceExt := filepath.Ext(destSourcePath)
		destSourcePath = strings.TrimSuffix(destSourcePath, sourceExt)

		// Open header and source files at dest
		destSourceFile, err := os.Create(destSourcePath + sourceExt)
		defer destSourceFile.Close()
		destHeaderFile, err := os.Create(destSourcePath + ".h")
		defer destHeaderFile.Close()

		if err != nil {
			log.Fatal(err)
		}

		_, destSourceFilename := filepath.Split(destSourcePath)

		// Write header guards and initial namespace data
		destHeaderFile.WriteString("#ifndef FILES_DATA_H\n#define FILES_DATA_H\n#include <string>\n")
		destHeaderFile.WriteString(fmt.Sprintf("\nnamespace %s { \n\tnamespace StaticResources {\n", namespace))
		defer destHeaderFile.WriteString("\n\t}\n}\n\n#endif\n")

		// Write initial source file data
		destSourceFile.WriteString(fmt.Sprintf("#include \"%s.h\"\nnamespace %s { \n\tnamespace StaticResources {\n", destSourceFilename, namespace))
		defer destSourceFile.WriteString("\n\t}\n}")

		innerPadding := "\t\t"

		// Iterate through all the files in the source and goblin them uup
		walkErr := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
			// Say no to directories guys
			if info.IsDir() {
				return nil
			}

			resBuf, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			path = filepath.ToSlash(path)

			destSourceFile.WriteString(fmt.Sprintf("%s// Data for file: %s\n", innerPadding, path))

			// TODO: Handle other symbols that are allowed in path names but not in variable declarations within C/C++
			pathAsSnake := strings.ReplaceAll(path, "/", "_")
			pathAsSnake = strings.ReplaceAll(pathAsSnake, ".", "_")
			pathAsSnake = strings.ReplaceAll(pathAsSnake, " ", "_")

			// Write declarations
			destHeaderFile.WriteString(fmt.Sprintf("%sextern const %s %s_path;\n", innerPadding, stringType, pathAsSnake))
			destSourceFile.WriteString(fmt.Sprintf("%sextern unsigned int %s_len;\n", innerPadding, pathAsSnake))
			destHeaderFile.WriteString(fmt.Sprintf("%sextern unsigned char %s[];\n", innerPadding, pathAsSnake))
			//extern const std::string Resources_default_layout_ini_path;
			//extern unsigned int Resources_default_layout_ini_len;
			//extern unsigned char Resources_default_layout_ini[];

			// Write resource section impl
			destSourceFile.WriteString(fmt.Sprintf("%sconst %s %s_path = R\"(%s)\";\n", innerPadding, stringType, pathAsSnake, path))
			destSourceFile.WriteString(fmt.Sprintf("%sunsigned int %s_len = %d;\n", innerPadding, pathAsSnake, len(resBuf)))
			destSourceFile.WriteString(fmt.Sprintf("%sunsigned char %s[] = {", innerPadding, pathAsSnake))

			// Now convert the file data to C style char[]
			var buf strings.Builder
			if len(resBuf) > 0 {
				buf.Grow(len(resBuf)*6 - 2)
				for i, b := range resBuf {
					if i%12 == 0 {
						buf.WriteString("\n" + innerPadding + "\t")
					}
					fmt.Fprintf(&buf, "0x%02x", b)
					if i < len(resBuf)-1 {
						buf.WriteString(",")
					}
				}
			}

			result := buf.String()
			destSourceFile.WriteString(result)

			destSourceFile.WriteString("\n" + innerPadding + "};\n\n")

			return nil
		})

		if walkErr != nil {
			log.Fatal(walkErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(restocsourceCmd)

	restocsourceCmd.Flags().StringP(restocsourceCmdInfo.FlagNameResDir, restocsourceCmdInfo.FlagNameResDirShort, "", "directory that contains files to be converted")
	restocsourceCmd.Flags().StringP(restocsourceCmdInfo.FlagNameDest, restocsourceCmdInfo.FlagNameDestShort, "", "destination file to write to. Can specify any extension, but the header file will be written to .h next to the specified destination file")
	restocsourceCmd.Flags().StringP(restocsourceCmdInfo.FlagNameNamespace, restocsourceCmdInfo.FlagNameNamespaceShort, "", "namespace of declarations")
	restocsourceCmd.Flags().StringP(restocsourceCmdInfo.FlagNameStringType, restocsourceCmdInfo.FlagNameStringTypeShort, "std::string", "the string type used in your codebase")

	restocsourceCmd.MarkFlagRequired(restocsourceCmdInfo.FlagNameResDir)
	restocsourceCmd.MarkFlagRequired(restocsourceCmdInfo.FlagNameDest)
	restocsourceCmd.MarkFlagRequired(restocsourceCmdInfo.FlagNameNamespace)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
