// Package main provides a CLI tool to remove comments from Go source code.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/nogocomments/internal/commentremover"
	"github.com/pierow2k/nogocomments/internal/filereader"
)

var (
	// Command-line flags for configuring application behavior.
	debugFlag   = flag.Bool("debug", false, "Enable debug logging level")
	fileFlag    = flag.String("file", "", "File path to read text from")
	pasteFlag   = flag.Bool("paste", false, "Read text from clipboard")
	versionFlag = flag.Bool("version", false, "Display version information")

	// ErrNoInputMethod is returned when no input method is specified.
	ErrNoInputMethod = errors.New("no input method specified")

	// loggingLevel manages the application's logging level dynamically.
	loggingLevel = new(slog.LevelVar)
)

// Application build-time metadata. These are set during the build process.
var (
	BuildDate     = "YYYY-MM-DDTHH:MM:SSZ"
	CopyrightDate = "YYYY"
	Version       = "X.X.X"
)

// initializeLoggingAndFlags configures command-line flags and initializes logging.
//
//nolint:errcheck
func initializeLoggingAndFlags() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Println("  --file <path>\tFile path to read text from")
		fmt.Println("  --paste\t\tRead text from clipboard")
		fmt.Println("  --debug\t\tEnable debug logging level")
		fmt.Println("  --version\t\tDisplay version information")
		fmt.Println("  --help\t\tShow usage information")
	}

	flag.Parse()

	// Configure logging with a text handler and default error level.
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel})
	slog.SetDefault(slog.New(h))

	loggingLevel.Set(slog.LevelError)

	// Enable debug logging if the debug flag is set.
	if *debugFlag {
		loggingLevel.Set(slog.LevelDebug)
	}
}

// displayVersion prints the application's version and build metadata.
func displayVersion() {
	fmt.Printf("nogocomments - Removes comments from Go code.\n\n")
	fmt.Printf("Version %s - Build Date %s\n", Version, BuildDate)
	fmt.Printf("Copyright (c) %s Pierow2k\n", CopyrightDate)
	fmt.Println("Distributed under the MIT License")
	os.Exit(0)
}

// displayUsage prints usage instructions and exits the application.
func displayUsage() {
	flag.Usage()
	os.Exit(1)
}

// readInputText reads input text from a file or clipboard based on user
// specified flags. Returns the input text or an error if reading fails.
func readInputText() (string, error) {
	if *fileFlag != "" {
		slog.Debug("reading text from file", "file", *fileFlag)

		fr := &filereader.RealFileReader{}

		text, err := filereader.ReadFile(fr, *fileFlag)
		if err != nil {
			return "", fmt.Errorf("failed to read from file '%s': %w", *fileFlag, err)
		}

		return text, nil
	}

	if *pasteFlag {
		slog.Debug("reading text from clipboard")

		text, err := clipboard.ReadAll()
		if err != nil {
			return "", fmt.Errorf("failed to read from clipboard: %w", err)
		}

		return text, nil
	}

	return "", ErrNoInputMethod
}

// processText removes comments from the given input text using the
// commentremover package. Returns the processed text or an error if the
// operation fails.
func processText(text string) (string, error) {
	processedText, err := commentremover.RemoveComments(text)
	if err != nil {
		return "", fmt.Errorf("failed to remove comments from source: %w", err)
	}

	return processedText, nil
}

// main is the entry point for the nogocomments application.
func main() {
	// Initialize logging and parse command-line flags.
	initializeLoggingAndFlags()

	// Display version information if the version flag is set.
	if *versionFlag {
		displayVersion()

		return
	}

	// Show usage instructions if no input method is provided.
	if !*pasteFlag && *fileFlag == "" {
		displayUsage()

		return
	}

	// Read input text from the specified source.
	text, err := readInputText()
	if err != nil {
		slog.Error("failed to read input", "error", err)
		os.Exit(1)
	}

	// Process the input text to remove comments.
	processedText, err := processText(text)
	if err != nil {
		slog.Error("failed to process text", "error", err)
		slog.Error("incomplete or non-go source code in input")
		os.Exit(1)
	}

	// Output the processed text to the console.
	fmt.Println(processedText)
	slog.Debug("application completed successfully")
}
