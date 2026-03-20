// Package filereader provides functionality for reading file contents.
// It uses dependency injection for filesystem operations and logging,
// allowing for easy testing and flexible configuration.
package filereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// fileSystem abstracts filesystem operations for reading files.
// This interface allows mocking filesystem access in tests.
type fileSystem interface {
	Open(name string) (io.ReadCloser, error)
}

// logger defines the interface for logging error messages.
type logger interface {
	Errorf(format string, args ...any)
}

// FileReader reads content from a file. It uses dependency injection
// for filesystem operations and logging, allowing for easy testing
// and flexible configuration.
type FileReader struct {
	fileSystem fileSystem
	filePath   string
	logger     logger
}

// New creates a new FileReader with the provided options.
// If no options are specified, defaults are used:
//   - fileSystem: osFS (real filesystem)
//   - logger: discardLogger (no logging)
//   - filePath: empty string (must be set via WithFile)
func New(opts ...Option) *FileReader {
	fileReader := &FileReader{
		fileSystem: osFS{},
		logger:     DiscardLogger{},
	}
	for _, opt := range opts {
		opt(fileReader)
	}

	return fileReader
}

// Option is a functional option for configuring a FileReader.
// Options can be passed to New to customize the FileReader's behavior.
type Option func(*FileReader)

// WithFile returns an Option that sets the file path to read.
// This is the primary way to specify which file the FileReader should read.
func WithFile(filePath string) Option {
	return func(f *FileReader) { f.filePath = filePath }
}

// WithFileSystem sets the filesystem for the FileReader.
func WithFileSystem(fileSystem fileSystem) Option {
	return func(f *FileReader) { f.fileSystem = fileSystem }
}

// WithLogger returns an Option that sets the logger for the FileReader.
// Use this to provide custom error logging behavior.
func WithLogger(l logger) Option {
	return func(f *FileReader) { f.logger = l }
}

// osFS implements fileSystem using the operating system's filesystem.
type osFS struct{}

// Open opens the named file using os.Open.
// Error wrapping is intentionally skipped to preserve os.Open's error types.
func (osFS) Open(name string) (io.ReadCloser, error) { return os.Open(name) } //nolint:wrapcheck

// DiscardLogger implements logger by discarding all log output.
// It is used as the default logger when none is provided via options.
type DiscardLogger struct{}

// Errorf implements logger.Errorf by doing nothing.
func (DiscardLogger) Errorf(_ string, _ ...any) {}

// ReadFile reads and returns the entire content of the configured file.
// It returns the file content as a string with line endings preserved,
// or an error if the file cannot be opened or read.
//
// The returned string includes newline characters after each line,
// including the last line.
func (fr *FileReader) ReadFile() (string, error) {
	file, err := fr.fileSystem.Open(fr.filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", fr.filePath, err)
	}

	defer func() {
		if c, ok := file.(io.Closer); ok {
			if err := c.Close(); err != nil {
				fr.logger.Errorf("failed to close file %s: %v\n", fr.filePath, err)
			}
		}
	}()

	var text strings.Builder

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s: %w", fr.filePath, err)
	}

	return text.String(), nil
}
