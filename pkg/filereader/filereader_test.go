package filereader_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pierow2k/nogocomments/pkg/filereader"
)

type mockFS struct {
	files map[string]io.ReadCloser
}

func (m *mockFS) Open(name string) (io.ReadCloser, error) {
	if f, ok := m.files[name]; ok {
		return f, nil
	}

	return nil, os.ErrNotExist
}

// TestFileReader_ReadFile provides unit tests for the ReadFile method of
// the FileReader.
//
//nolint:funlen
func TestFileReader_ReadFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		filePath  string
		mockFiles map[string]string
		want      string
		wantErr   bool
	}{
		{
			name:     "read existing file",
			filePath: "/path/to/file.txt",
			mockFiles: map[string]string{
				"/path/to/file.txt": "Hello, World!",
			},
			want:    "Hello, World!\n",
			wantErr: false,
		},
		{
			name:     "read file with multiple lines",
			filePath: "multiline.txt",
			mockFiles: map[string]string{
				"multiline.txt": "Line 1\nLine 2\nLine 3",
			},
			want:    "Line 1\nLine 2\nLine 3\n",
			wantErr: false,
		},
		{
			name:      "file does not exist",
			filePath:  "nonexistent.txt",
			mockFiles: map[string]string{},
			wantErr:   true,
		},
		{
			name:     "read empty file",
			filePath: "empty.txt",
			mockFiles: map[string]string{
				"empty.txt": "",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			files := make(map[string]io.ReadCloser)
			for name, content := range testCase.mockFiles {
				files[name] = io.NopCloser(strings.NewReader(content))
			}

			//nolint:varnamelen
			fr := filereader.New(
				filereader.WithFileSystem(&mockFS{files: files}),
				filereader.WithFile(testCase.filePath),
				filereader.WithLogger(filereader.DiscardLogger{}),
			)

			got, gotErr := fr.ReadFile()
			if (gotErr != nil) != testCase.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", gotErr, testCase.wantErr)

				return
			}

			if got != testCase.want {
				t.Errorf("ReadFile() = %q, want %q", got, testCase.want)
			}
		})
	}
}
