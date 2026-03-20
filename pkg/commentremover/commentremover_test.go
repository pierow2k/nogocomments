package commentremover_test

import (
	"strings"
	"testing"

	"github.com/pierow2k/nogocomments/pkg/commentremover"
)

// TestRemoveComments provides unit tests for the RemoveComments function.
//
//nolint:funlen
func TestRemoveComments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		commentMarker string // Use this for tests where comment removal is verified
		want          string // Use this for direct output comparison when applicable
		wantErr       bool
	}{
		{
			name: "remove single line comments",
			input: `package main

// This is a uniqueCommentMarker single line comment
func main() {
    fmt.Println("Hello, world!")
}`,
			commentMarker: "uniqueCommentMarker", // This should not appear in the output
			wantErr:       false,
		},
		{
			name: "remove multi-line comments",
			input: `package main

/*
This is a multi-line comment containing uniqueMultiLineMarker.
*/
func main() {
    fmt.Println("Hello, world!")
}`,
			commentMarker: "uniqueMultiLineMarker", // This should not appear in the output
			wantErr:       false,
		},
		{
			name: "input without package declaration",
			input: `// commentWithUniqueNoPackageMarker
func example() {
    fmt.Println("Example function")
}`,
			commentMarker: "commentWithUniqueNoPackageMarker", // This should not appear in the output
			wantErr:       false,
		},
		{
			name:    "invalid Go code",
			input:   `package main func main() {`,
			wantErr: true, // Expect an error for invalid Go code
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := commentremover.RemoveComments(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("RemoveComments() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if testCase.commentMarker != "" && strings.Contains(got, testCase.commentMarker) {
				t.Errorf("RemoveComments() result should not contain comment marker '%v',got = %v",
					testCase.commentMarker, got)
			}

			if testCase.want != "" && got != testCase.want {
				t.Errorf("RemoveComments() got = %v, want %v", got, testCase.want)
			}
		})
	}
}
