package filereader_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pierow2k/nogocomments/internal/filereader"
	"github.com/stretchr/testify/assert"
)

// ExampleReadFile demonstrates the use of the ReadFile function from the
// filereader package.
func ExampleReadFile() {
	// Setup: Typically, you wouldn't use a real file in examples or tests,
	// but for simplicity, assume there's a file named "example.txt" with
	// some content. In a real-world scenario, you'd use a mock or temp
	// file.
	//
	// Create a RealFileReader instance.
	fr := &filereader.RealFileReader{}

	// Assuming "example.txt" is a file in your current directory with
	// some text.
	text, err := filereader.ReadFile(fr, "example.txt")
	if err != nil {
		fmt.Println("Failed to read file:", err)

		return
	}

	fmt.Print(text)
	// Output: test data
}

// mockFileReader is a mock implementation of FileReader for testing purposes.
type mockFileReader struct {
	content string
	err     error
}

func (mfr *mockFileReader) Open(_ string) (io.Reader, error) {
	if mfr.err != nil {
		return nil, mfr.err
	}

	return strings.NewReader(mfr.content), nil
}

//nolint:funlen,err113
func TestReadFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		err     error
		want    string
		wantErr bool
	}{
		{"Basic test", "Hello, World!\n", nil, "Hello, World!\n", false},
		{"Empty file", "", nil, "", false},
		{"File read error", "", errors.New("mock read error"), "", true},
		{"Permission denied", "", os.ErrPermission, "", true},
		{"Ends with newline", "Line 1\nLine 2\nLine 3\n", nil, "Line 1\nLine 2\nLine 3\n", false},
		{"Does not end with newline", "Line 1\nLine 2\nLine 3", nil, "Line 1\nLine 2\nLine 3\n", false},
		{
			name:    "Overlong UTF-8 encoding",
			content: "Valid UTF-8 text\xC0\x81Overlong UTF-8 encoding",
			err:     nil,
			want:    "Valid UTF-8 text\xC0\x81Overlong UTF-8 encoding\n",
			wantErr: false,
		},
		{
			name:    "Incomplete multi-byte characters",
			content: "Valid UTF-8 text\xE0\xA0Incomplete multi-byte characters",
			err:     nil,
			want:    "Valid UTF-8 text\xE0\xA0Incomplete multi-byte characters\n",
			wantErr: false,
		},
		{
			name:    "Invalid continuation byte",
			content: "Valid UTF-8 start\xC2\x50Invalid continuation",
			err:     nil,
			want:    "Valid UTF-8 start\xC2\x50Invalid continuation\n",
			wantErr: false,
		},
		{
			name:    "Invalid UTF-8 sequence",
			content: "Valid text\xE3\x80\x22Invalid sequence",
			err:     nil,
			want:    "Valid text\xE3\x80\x22Invalid sequence\n",
			wantErr: false,
		},
	}

	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			mfr := &mockFileReader{content: testTable.content, err: testTable.err}
			got, err := filereader.ReadFile(mfr, "dummy/path.txt")

			if (err != nil) != testTable.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, testTable.wantErr)

				return
			}

			assert.Equal(t, testTable.want, got)
		})
	}
}

func TestReadFile_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	mockFR := &mockFileReader{
		err: os.ErrNotExist,
	}
	_, err := filereader.ReadFile(mockFR, "nonexistent_file.txt")
	assert.ErrorIs(t, err, os.ErrNotExist)
}
