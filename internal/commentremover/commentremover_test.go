package commentremover_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pierow2k/nogocomments/internal/commentremover"
)

func ExampleRemoveComments() {
	sourceCode := `package main
// This is a single-line comment.
func main() {
	// Another comment.
	fmt.Println("Hello, World!")
	/* A multi-line
	comment */
}`

	// Remove comments from the source code.
	cleanSource, err := commentremover.RemoveComments(sourceCode)
	if err != nil {
		fmt.Println("Error removing comments:", err)

		return
	}

	fmt.Println(cleanSource)
	// Output:
	// package main
	//
	// func main() {
	//
	// 	fmt.Println("Hello, World!")
	//
	// }
}

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

	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got, err := commentremover.RemoveComments(testTable.input)
			if (err != nil) != testTable.wantErr {
				t.Errorf("RemoveComments() error = %v, wantErr %v", err, testTable.wantErr)

				return
			}

			if testTable.commentMarker != "" && strings.Contains(got, testTable.commentMarker) {
				t.Errorf("RemoveComments() result should not contain comment marker '%v', got = %v", testTable.commentMarker, got)
			}

			if testTable.want != "" && got != testTable.want {
				t.Errorf("RemoveComments() got = %v, want %v", got, testTable.want)
			}
		})
	}
}
