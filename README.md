# goblin
A small, cross-platform executable that contains a collection of small commands under one executable. Designed to make seemingly simple tasks easier.

```
go install github.com/WiggleWizard/goblin
```

## Current Commands
## echo
Prints a message, without a newline and without quotes
```
Why do we need an echo command when every OS has one already? Cross-platformness is the answer,
along with some weirdness on Windows. Windows can goblin deez nutz.

Usage:
  goblin echo [flags] [message]

Flags:
  -f, --tofile string   regex expression to use as a needle
```
## genusage
Generates a total list of all commands, usages and descriptions
```
Usage:
  goblin genusage [flags]

Flags:
  -h, --help   help for genusage
      --md     generate each command usage inside a markdown code block
```
## restocsource
Generates a C/C++ source and header that contains embedded resource data
```
Usage:
  goblin restocsource [flags]

Flags:
  -d, --destination string   destination file to write to. Can specify any extension, but the header file will be written to .h next to the specified destination file
  -n, --namespace string     namespace of declarations
  -s, --sourcedir string     directory that contains files to be converted
  -t, --stringtype string    the string type used in your codebase (default "std::string")
```
## rm
remove files or directories
```
rm removes each specified file. By default, it does not remove directories.

Usage:
  goblin rm [flags]

Flags:
  -f, --force       ignore nonexistent files and arguments, never prompt
  -r, --recursive   remove directories and their contents recursively
```
## sub
Regex replace in strings or files
```
sub is a small utility that uses regex and regex groups to substitute
text inside either a piped in string or a file with arbitrarily numbered inputs.
The results are printed via stdout.

The regex flavor is Go, which contains group names. In the context of sub
this allows specifying exact replacement strings within the CLI args by placing
the strings at the correct CLI arg index. Here's a simple example:

	> goblin echo Hello World | goblin sub -e "Hello (?P<1>(?s).*)" "Foo"
	Hello Foo

	> goblin echo Hello World | goblin sub -e "(?P<2>Hello) (?P<1>World)" "Foo" "Bar"
	Bar Foo

sub groups are labelled sequentially with the [...substitutes] passed with
group 0 specifically referring to stdin/piped data:

	> goblin echo Hello World | goblin sub -e "Hello(?P<0>\s)World"
	HelloHello WorldWorld

File data can also be loaded as a haystack by using the --input flag. If this flag is used then
stdin is not operated on as a haystack:

	README.md:
		Goblin is pretty cool bro.

	> goblin echo deez nutz | goblin sub -i README.md -e "(?P<0>is pretty cool) -o README.md"

	README.md:
		Goblin deez nutz bro.


Usage:
  goblin sub [flags] [...substitutes]

Flags:
  -e, --expression string   regex expression to use as a needle
  -i, --input string        use file contents as haystack. Not specifying this flag will cause stdin to be the
                            haystack and regex group <0> to be available as a replacement group
  -o, --output string       output file to write the results to
```
## time
Prints the current date and time to stdout
```
Usage:
  goblin time [flags]

Flags:
  -f, --format string   see https://man7.org/linux/man-pages/man3/strftime.3.html for formatting
  -u, --unix            outputs the date and time in Unix timestamp
```
