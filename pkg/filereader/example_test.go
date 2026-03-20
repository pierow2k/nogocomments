//nolint:testableexamples
package filereader_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/pierow2k/nogocomments/pkg/filereader"
)

// New creates a new FileReader instance. It is most commonly used in
// combination with the WithFile option.
func ExampleNew() {
	fr := filereader.New()
	_ = fr
}

// WithFile sets the file path to read. This is the primary way to specify
// which file the FileReader should read.
func ExampleWithFile() {
	fr := filereader.New(
		filereader.WithFile("/path/to/file.txt"),
	)
	_ = fr
}

// WithFileSystem sets the filesystem that is used by the FileReader. If
// no fileSystem is specified, New defaults to using the operating system
// filesystem. A mock filesystem is used for demonstration purposes in this
// example.
func ExampleWithFileSystem() {
	files := make(map[string]io.ReadCloser)
	fr := filereader.New(
		filereader.WithFileSystem(&mockFS{files: files}),
	)
	_ = fr
}

// WithLogger sets the logger by the FileReader. New defaults to using an
// internal DiscardLogger that discards all log output. In this example,
// the DiscardLogger is used but any suitable alternate logger can be
// specified.
func ExampleWithLogger() {
	fr := filereader.New(
		filereader.WithLogger(filereader.DiscardLogger{}),
	)
	_ = fr
}

// ReadFile reads and returns the entire content of a file. A mock
// filesystem is used for demonstration purposes in this example.
func ExampleFileReader_ReadFile() {
	filePath := "/path/to/file.txt"
	content := "Hello, World!"
	files := map[string]io.ReadCloser{
		filePath: io.NopCloser(strings.NewReader(content)),
	}

	fr := filereader.New(
		filereader.WithFile(filePath),
		filereader.WithFileSystem(&mockFS{files: files}),
	)

	fileContents, err := fr.ReadFile()
	if err != nil {
		panic(err)
	}

	fmt.Println(fileContents)

	// Output: Hello, World!
}
