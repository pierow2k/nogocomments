// Package commentremover removes comments from Go source code. It handles
// both complete packages and individual code snippets using the go/parser,
// go/ast, go/token, and go/printer packages.
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

const dummyPackage = "package main\n"

// ensurePackageDeclaration prepends a "package main" declaration to
// sourceCode if it doesn't already start with one. This is necessary to
// properly parse a snippet since go/parser requires a package declaration.
// It returns the potentially modified source and a boolean indicating
// whether the dummy package was added.
func ensurePackageDeclaration(sourceCode string) (string, bool) {
	var sourceScanner scanner.Scanner

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(sourceCode))
	sourceScanner.Init(file, []byte(sourceCode), nil, scanner.ScanComments)

	// Scan through tokens to find the first non-comment token.
	for {
		_, tok, _ := sourceScanner.Scan()
		if tok == token.EOF {
			break
		}
		// Return immediately once we determine if package exists.
		if tok != token.COMMENT && tok != token.PACKAGE {
			return dummyPackage + sourceCode, true
		}
		// If we find a package token, no need to prepend dummy package.
		if tok == token.PACKAGE {
			return sourceCode, false
		}
	}

	// Default to not modifying source code if it's unclear what action to
	// take. This might happen for files that only contain comments.

	return sourceCode, false
}

// parseSourceCode parses sourceCode into an AST using the provided file set.
func parseSourceCode(fset *token.FileSet, sourceCode string) (*ast.File, error) {
	file, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing source code: %w", err)
	}

	return file, nil
}

// removeCommentsFromAST removes all comment groups from file in-place.
func removeCommentsFromAST(file *ast.File) {
	file.Comments = []*ast.CommentGroup{}
}

// formatAST converts the AST back into a Go source code string.
func formatAST(file *ast.File, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, file); err != nil {
		return "", fmt.Errorf("error formatting source code: %w", err)
	}

	return buf.String(), nil
}

// removeDummyPackage removes the leading "package main\n" from sourceCode.
func removeDummyPackage(sourceCode string) string {
	return strings.TrimPrefix(sourceCode, dummyPackage)
}

// RemoveComments removes all comments from the provided Go source code.
// It handles both complete packages and standalone code snippets. If the
// source lacks a package declaration, a temporary one is added for parsing
// and removed from the output.
func RemoveComments(sourceCode string) (string, error) {
	fset := token.NewFileSet()
	sourceCode, prefixed := ensurePackageDeclaration(sourceCode)

	file, err := parseSourceCode(fset, sourceCode)
	if err != nil {
		return "", err
	}

	removeCommentsFromAST(file)

	result, err := formatAST(file, fset)
	if err != nil {
		return "", err
	}

	if prefixed {
		result = removeDummyPackage(result)
	}

	return result, nil
}
