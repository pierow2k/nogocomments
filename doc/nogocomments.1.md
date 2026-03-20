% nogocomments(1) Version {{version}}| General Commands Manual <!-- markdownlint-disable MD041 -->
%
% {{version}}

# NAME

**nogocomments** - Instantly Remove Comments from Your Go Code

# SYNOPSIS

| **nogocomments** \[**INPUT_FILE**\] \[**OPTIONS**\]

# DESCRIPTION

`nogocomments` removes comments from Go source code. It reads Go code from
a file or the system clipboard and writes the result to standard output. It
supports both complete packages and standalone code snippets.

# OPTIONS

| Short | Long        | Description                              |
| :---: | :---------- | :--------------------------------------- |
| `-d`  | `--debug`   | Enable debug (verbose) logging           |
| `-h`  | `--help`    | Show help                                |
| `-p`  | `--paste`   | Read code from clipboard                 |
| `-v`  | `--version` | Show version, build details, and license |

# EXAMPLES

**Remove comments from Go source code that has been copied to the clipboard:**

`nogocomments --paste`

**Remove comments from a Go file and print the result to the terminal:**

`nogocomments --file /path/to/your/file.go`

**Remove comments from a Go file and write the result to a new file:**

`nogocomments --file /path/to/your/source.go > newfile.go`

**Enable debug mode for more detailed logs:**

`nogocomments --debug --file /path/to/your/file.go`

**Print the version, build details, and license information:**

`nogocomments --version`

# BUGS

Report issues at the project tracker: https://github.com/pierow2k/nogocomments/issues

# COPYRIGHT

Copyright (C) {{copyright_date}} Pierow2K.

# LICENSE

nogocomments is distributed under the MIT License.  
See the `LICENSE` file for details.
