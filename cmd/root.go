// Package cmd implements the CLI for nogocomments, a tool that removes
// comments from Go source code. Call Execute to run the application.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/nogocomments/pkg/commentremover"
	"github.com/spf13/cobra"
)

// Configuration stores the configuration parsed from command-line flags.
type Configuration struct {
	filePath     string // filePath is the path to the Go source file to process.
	useClipboard bool   // useClipboard indicates whether to read input from the clipboard.
}

var (
	cfg Configuration // cfg is the global Configuration instance.

	// errMutuallyExclusive is returned when both a file path and the paste
	// flag are specified simultaneously.
	errMutuallyExclusive = errors.New("paste and file are mutually exclusive")

	// errNoInputMethod is returned when neither a file path nor the paste
	// flag is specified.
	errNoInputMethod = errors.New("no input method specified")
)

// BuildDate, CopyrightDate, Version, and License contain build information.
// These are placeholder values that are overwritten by linker flags at
// compile time.
var (
	BuildDate     = "YYYY-MM-DDTHH:MM:SS-0000"
	CopyrightDate = "2026"
	Version       = "v0.0.0-dev"
	License       = "Licensed under the MIT License <https://opensource.org/licenses/MIT>"
)

// rootCmd represents the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "nogocomments [INPUT_FILE]",
	Short: "Remove comments from Go source code.",
	Long: `nogocomments removes comments from Go source code.
It reads Go code from a file or the system clipboard and writes the
result to standard output. It supports both complete packages and
standalone code snippets.`,
	Example: `  # Remove comments from a file
  nogocomments somecode.go

  # Remove comments from code on the clipboard
  nogocomments --paste`,
	Version: fmt.Sprintf(
		"%s - built %s\nCopyright © %s Pierow2k\n%s",
		Version, BuildDate, CopyrightDate, License,
	),
	Args: cobra.MaximumNArgs(1),
	RunE: runFunction,
}

// Execute is the entry point for the CLI. It processes command-line
// arguments and exits with a non-zero status code on error.
func Execute() {
	rootCmd.InitDefaultHelpFlag()
	rootCmd.Flags().Lookup("help").Usage = "Show help"
	rootCmd.InitDefaultVersionFlag()
	rootCmd.Flags().Lookup("version").Usage = "Show version, build details, and license"

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// init registers the command-line flags for the root command.
func init() {
	rootCmd.Flags().BoolVarP(&cfg.useClipboard, "paste", "p", false, "Read code from the system clipboard")
}

// runFunction implements the root command. It reads Go source code from
// a file or the clipboard, removes comments, and writes the result to
// standard output.
//
// Errors are returned in the following cases:
//   - Both a file path and the paste flag are specified (mutually exclusive)
//   - Neither a file path nor the paste flag is specified
//   - Reading from the file or clipboard fails
//   - Comment removal fails
func runFunction(_ *cobra.Command, args []string) error {
	var (
		sourceCode string
		err        error
	)

	if len(args) > 0 {
		cfg.filePath = args[0]
	}

	switch {
	case cfg.useClipboard && cfg.filePath != "":
		return errMutuallyExclusive
	case !cfg.useClipboard && cfg.filePath == "":
		return errNoInputMethod
	}

	if cfg.useClipboard {
		sourceCode, err = clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read from clipboard: %w", err)
		}
	} else {
		fileContent, err := os.ReadFile(cfg.filePath)
		if err != nil {
			return fmt.Errorf("file read failed: %w", err)
		}

		sourceCode = string(fileContent)
	}

	result, err := commentremover.RemoveComments(sourceCode)
	if err != nil {
		return fmt.Errorf("failed to remove comments from source: %w", err)
	}

	_, _ = fmt.Fprintln(os.Stdout, result)

	return nil
}
