// Package commentremover provides a utility for removing comments from
// Go source code. It supports handling both complete Go packages and
// individual snippets of Go code. The primary functionality is offered
// through the RemoveComments function, which takes a string of Go source
// code and returns a version of that string with all comments removed.
// This process involves temporarily adding a dummy package declaration
// if necessary to facilitate parsing, converting the source code into an
// Abstract Syntax Tree (AST), removing comment nodes from the AST, and
// then formatting the AST back into a source code string.
//
// The package leverages the Go standard library's "go/parser", "go/token",
// "go/ast", and "go/printer" packages to parse source code, manipulate
// the AST, and format the AST back into source code. This makes it
// particularly useful for preprocessing Go source code for further
// analysis, transformation, or any other context where comments are not
// needed.
//
// Functions:
//   - RemoveComments: The main function that removes comments from the
//     provided Go source code. It handles adding and removing a dummy
//     package declaration as necessary to ensure the code can be parsed
//     and processed correctly. Returns the modified source code without
//     comments or an error if parsing or processing fails.
//   - checkAndPrefixSource: A helper function that checks for the presence
//     of a package declaration at the beginning of the source code. If
//     absent, it prefixes the source code with a dummy "package main"
//     declaration.
//   - parseSourceCode: Parses the provided Go source code into an AST, using
//     the "go/parser" package. Returns the parsed AST or an error if
//     parsing fails.
//   - removeCommentsFromAST: Removes all comment nodes from the provided
//     AST, modifying the AST in place.
//   - formatAST: Formats the modified AST back into a Go source code string,
//     using the "go/printer" package. Returns the formatted source code or
//     an error if formatting fails.
//   - removeDummyPackage: Removes the dummy package declaration added by
//     checkAndPrefixSource, if it was added. This ensures the output
//     matches the original structure of the input source code, minus the
//     comments.
//
// This package is useful for developers needing to preprocess Go source
// code to remove comments, whether for analysis, transformation, or other
// purposes where comments are not necessary.
package commentremover

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"strings"
)

// RemoveComments takes a string containing Go source code and removes all
// comments from it. This function can handle both complete packages and
// individual snippets of Go code. If the source code does not start with a
// package declaration, a dummy "package main" declaration is temporarily
// added to facilitate parsing. This dummy declaration is removed from the
// output if it was added. The process involves converting the source code
// into an Abstract Syntax Tree (AST), removing comment nodes from the AST,
// and then formatting the modified AST back into a source code string.
//
// The function is designed for situations where comments are not needed,
// such as preprocessing source code for analysis or transformation. It
// provides a convenient way to clean up source code by stripping out
// comments without altering the structure or functionality of the code
// itself.
//
// Parameters:
//   - sourceCode: A string containing the Go source code from which comments
//     are to be removed.
//
// Returns:
//   - A string of the modified source code with all comments removed.
//   - An error if the source code could not be parsed or processed for any
//     reason.
//
// Note: The function returns an error if the provided source code is
// invalid or if any issues are encountered during the parsing, AST
// manipulation, or formatting stages of the process.
func RemoveComments(sourceCode string) (string, error) {
	fset := token.NewFileSet()
	sourceCode, prefixed := checkAndPrefixSource(sourceCode)

	file, err := parseSourceCode(fset, sourceCode)
	if err != nil {
		return "", err
	}

	removeCommentsFromAST(file)

	result, err := formatAST(&file, fset)
	if err != nil {
		return "", err
	}

	if prefixed {
		result = removeDummyPackage(result)
	}

	return result, nil
}

// checkAndPrefixSource examines the provided source code to determine
// whether it starts with a package declaration or leading comments.
// If the source code does not begin with a package declaration, a dummy
// "package main" declaration is prepended to facilitate parsing. This
// is particularly useful for handling code snippets that might be valid
// Go code but do not conform to the expected file structure (e.g., a
// function with leading comments but without a package declaration).
// The function uses a scanner to identify the first significant token
// (ignoring comments) to decide whether to prepend the dummy package
// declaration. It returns the potentially modified source code and a
// boolean indicating whether the dummy package declaration was added.
//
// Parameters:
// - sourceCode: A string containing the Go source code to be processed.
//
// Returns:
//   - The modified source code with a dummy package declaration prepended if
//     necessary.
//   - A boolean value indicating whether the dummy package declaration was
//     added (true if added, false otherwise).
//
// Note: The function is designed to improve parsing compatibility with
// standalone Go code snippets by ensuring they have a valid package
// declaration. This allows for the flexible handling of various Go code
// structures, including those that start with comments or function
// declarations without a preceding package declaration.
func checkAndPrefixSource(sourceCode string) (string, bool) {
	const dummyPackage = "package main\n"

	var sourceScanner scanner.Scanner

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(sourceCode))
	sourceScanner.Init(file, []byte(sourceCode), nil, scanner.ScanComments)

	// Scan through tokens to find the first non-comment token
	for {
		_, tok, _ := sourceScanner.Scan()
		if tok == token.EOF {
			break
		}
		// If the first non-comment, non-package token is encountered, prepend dummy package
		if tok != token.COMMENT && tok != token.PACKAGE {
			return dummyPackage + sourceCode, true
		}
		// If we find a package token, no need to prepend dummy package
		if tok == token.PACKAGE {
			return sourceCode, false
		}
	}

	// Default to not modifying source code if it's unclear what action to
	// take. This might happen for files that only contain comments.

	return sourceCode, false
}

// parseSourceCode parses the given Go source code string into an AST
// (Abstract Syntax Tree) using a provided file set. It returns the parsed
// file as an *ast.File and any error encountered during parsing.
func parseSourceCode(fset *token.FileSet, sourceCode string) (*ast.File, error) {
	file, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing source code: %w", err)
	}

	return file, nil
}

// removeCommentsFromAST takes an *ast.File and removes all comment groups
// from it. This operation directly modifies the AST in-place, effectively
// stripping out all comments from the source code represented by the AST.
func removeCommentsFromAST(file *ast.File) {
	file.Comments = []*ast.CommentGroup{}
}

// formatAST takes a pointer to an *ast.File and a *token.FileSet,
// converting the AST back into a Go source code string. It returns the
// formatted source code as a string and any error encountered during
// formatting.
func formatAST(file **ast.File, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, *file); err != nil {
		return "", fmt.Errorf("error formatting source code: %w", err)
	}

	return buf.String(), nil
}

// removeDummyPackage removes the dummy package declaration ("package
// main\n") from the beginning of the source code string, if present.
// This step is necessary to revert the preprocessing done by
// checkAndPrefixSource when the source code did not originally include a
// package declaration. It ensures the output closely matches the input
// structure, minus the comments.
func removeDummyPackage(sourceCode string) string {
	const dummyPackage = "package main\n"

	return strings.Replace(sourceCode, dummyPackage, "", 1)
}
