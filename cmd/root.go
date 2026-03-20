// Package cmd implements the command-line interface for the nogocomments application.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/nogocomments/pkg/commentremover"
	"github.com/pierow2k/nogocomments/pkg/filereader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// noGoComments stores the configuration from command-line flags, the
// source code, and the processed source code.
type noGoComments struct {
	debug         bool   // Enable debug logging.
	file          string // File path to read input from.
	paste         bool   // Read input from the clipboard.
	sourceCode    string // Source code to remove comments from.
	noCommentCode string // Processed source code with comments removed.
}

var (
	// nogo is the global noGoComments instance for the application.
	nogo noGoComments
	// ErrMutuallyExclusive is returned when a file path is specified and the
	// paste flag is used.
	ErrMutuallyExclusive = errors.New("paste and file are mutually exclusive")
	// ErrNoInputMethod is returned when neither a file path nor paste
	// are specified.
	ErrNoInputMethod = errors.New("no input method specified")
)

// Set temporary values for versioning. These are overwritten by linker flags
// in the Makefile at compile time.
var (
	BuildDate     = "YYYY-MM-DDTHH:MM:SS-0000"
	CopyrightDate = "2026"
	Version       = "v0.0.0-dev"
	License       = "Licensed under the MIT License <https://opensource.org/licenses/MIT>"
)

// rootCmd is the main command for the nogocomments application.
// It represents the base command when called without any subcommands.
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
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if nogo.debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Debuging enabled")
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	},
	Args: cobra.MaximumNArgs(1), // Ensure exactly one argument is provided
	RunE: runFunction,
}

// Execute executes the root command and will exit with a non-zero status
// code if an error occurs. Execute adds any child commands to the root
// command automatically.
func Execute() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		// DisableLevelTruncation: true,
		FullTimestamp: true,
		PadLevelText:  true,
	})

	// Initialize the default help flag and modify its description
	rootCmd.InitDefaultHelpFlag()

	helpFlag := rootCmd.Flags().Lookup("help")
	if helpFlag != nil {
		helpFlag.Usage = "Show help"
	}

	// Initialize the default help flag and modify its description
	rootCmd.InitDefaultVersionFlag()

	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag != nil {
		versionFlag.Usage = "Show version, build details, and license"
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init registers command-line flags.
func init() {
	rootCmd.Flags().BoolVarP(&nogo.debug, "debug", "d", false, "Enable debug (verbose) logging")
	rootCmd.Flags().BoolVarP(&nogo.paste, "paste", "p", false, "Read code from clipboard")
}

// readInputFile reads source code from a file.
func readInputFile() error {
	logrus.Debugf("reading text from file: %s", nogo.file)

	fr := filereader.New(
		filereader.WithFile(nogo.file),
		filereader.WithLogger(logrus.StandardLogger()))

	sourceCode, err := fr.ReadFile()
	if err != nil {
		return fmt.Errorf("failed to read from file '%s': %w", nogo.file, err)
	}

	nogo.sourceCode = sourceCode

	return nil
}

// readInputClipboard reads source code from the clipboard.
func readInputClipboard() error {
	logrus.Debugf("reading text from clipboard")

	sourceCode, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read from clipboard: %w", err)
	}

	nogo.sourceCode = sourceCode

	return nil
}

// processSourceCode removes comments from Go source code using
// the commentremover package.
func processSourceCode() error {
	noCommentCode, err := commentremover.RemoveComments(nogo.sourceCode)
	if err != nil {
		return fmt.Errorf("failed to remove comments from source: %w", err)
	}

	nogo.noCommentCode = noCommentCode

	return nil
}

// runFunction is the function that is executed when the root command is called.
func runFunction(_ *cobra.Command, args []string) error {
	var err error

	if len(args) > 0 {
		nogo.file = args[0]
	}

	if nogo.paste && nogo.file != "" {
		return ErrMutuallyExclusive
	} else if !nogo.paste && nogo.file == "" {
		return ErrNoInputMethod
	}

	if nogo.paste {
		err = readInputClipboard()
	} else {
		err = readInputFile()
	}

	if err != nil {
		return err
	}

	err = processSourceCode()
	if err != nil {
		return err
	}

	fmt.Println(nogo.noCommentCode)
	logrus.Debug("comment removal completed")

	return nil
}
