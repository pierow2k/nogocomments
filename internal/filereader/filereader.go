// Package filereader provides abstractions for reading files,
// enabling easier testing and flexible file reading strategies.
package filereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// FileReader defines an interface for reading files.
type FileReader interface {
	Open(filePath string) (io.Reader, error)
}

// RealFileReader is a FileReader implementation using the OS package.
type RealFileReader struct{}

// Open opens a file at the given path and returns an io.Reader.
func (rf *RealFileReader) Open(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return file, nil
}

// ReadFile reads the content of a file using the provided FileReader.
// Returns the file content as a string or an error.
func ReadFile(fr FileReader, filePath string) (string, error) {
	file, err := fr.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	defer func() {
		if c, ok := file.(io.Closer); ok {
			if err := c.Close(); err != nil {
				// Log or handle the error
				fmt.Printf("failed to close file %s: %v\n", filePath, err)
			}
		}
	}()

	var text strings.Builder

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return text.String(), nil
}
