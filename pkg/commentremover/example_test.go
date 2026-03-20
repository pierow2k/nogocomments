package commentremover_test

import (
	"fmt"

	"github.com/pierow2k/nogocomments/pkg/commentremover"
)

// The RemoveComments function removes all comments from the provided Go
// source code.
func ExampleRemoveComments() {
	sourceCode := `// Package main provides an example of the classic Hello World program.
package main

// main is the entry point of the program and prints "Hello, World!" to the console.
func main() {
	/* We print the text "Hello, World!", but the text
	can be changed to print any message. */
	fmt.Println("Hello, World!")
}`

	cleanSource, err := commentremover.RemoveComments(sourceCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(cleanSource)
	// Output:
	// package main
	//
	// func main() {
	//
	// 	fmt.Println("Hello, World!")
	// }
}

// The RemoveComments function can also be used with code snippets.
func ExampleRemoveComments_snippet() {
	sourceCode := `func main() {
	// We print the text "Hello, World!", but the text
	// can be changed to print any message.
	fmt.Println("Hello, World!")
}`

	// Remove comments from the source code.
	cleanSource, err := commentremover.RemoveComments(sourceCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(cleanSource)
	// Output:
	// func main() {
	//
	// 	fmt.Println("Hello, World!")
	// }
}
