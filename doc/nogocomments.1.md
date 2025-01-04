% NOGOCOMMENTS(1) Version {{version}} | User Commands
%
% {{date}}

NAME
====

**nogocomments** - A utility tool for removing Go comments from source code.


SYNOPSIS
========

| **nogocomments** *\[OPTION\]*

DESCRIPTION
===========

nogocomments is a utility designed to clean up Go source code by removing
comments. It can process Go source code either from a specified file or
directly from the clipboard, effectively eliminating both line comments
(beginning with "//") and inline comments (whitespace followed by "//").
nogocomments writes its output to standard out, which can be redirected to
a file if necessary. The tool aims to provide a straightforward solution
for developers looking to clean up their code for documentation or sharing
purposes without comments.

OPTIONS
=======

**\-\-debug**
:   Enable debug level logging. This flag increases the verbosity of
    the logging output, useful for debugging issues or understanding
    the tool's operations in more detail.

**\-\-file** *\<file_path\>*
:   Specify the file path to read the Go source code from.

**\-\-paste**
:   Read Go source code from the system clipboard.

**\-\-version**
:   Display the version information of the nogocomments tool, including
    the build date.


EXAMPLES
========

Remove comments from Go source code copied to the clipboard:

```bash
nogocomments --paste
```

Process comments from a Go source file and output to the terminal:

```bash
nogocomments -file /path/to/your/file.go
```

Write processed output to a new file:

```bash
nogocomments -file /path/to/your/source.go > newfile.go
```

Enable debug mode for detailed logs:

```bash
nogocomments -debug -file /path/to/your/file.go
```

REPORTING BUGS
==============

Report bugs or suggest improvements via GitHub Issues:
[https://github.com/pierow2k/nogocomments/issues/](https://github.com/pierow2k/nogocomments/issues/)


COPYRIGHT
=========

Copyright (c) {{copyright_date}} Pierow2K. Released under the MIT License.
