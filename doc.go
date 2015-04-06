/*
Command pigeon generates Go parsers from a PEG grammar.

From Wikipedia [0]:

	A parsing expression grammar is a type of analytic formal grammar, i.e.
	it describes a formal language in terms of a set of rules for recognizing
	strings in the language.

Its features and syntax are inspired by the PEG.js project [1], while
the implementation is loosely based on [2].

	[0]: http://en.wikipedia.org/wiki/Parsing_expression_grammar
	[1]: http://pegjs.org/
	[2]: http://www.codeproject.com/Articles/29713/Parsing-Expression-Grammar-Support-for-C-Part

Command-line usage

The pigeon tool must be called with a PEG grammar file as defined
by the accepted PEG syntax below. The grammar may be provided by a
file or read from stdin. The generated parser is written to stdout
by default.

	pigeon [options] [GRAMMAR_FILE]

The following options can be specified:

	-debug : boolean, print debugging info to stdout (default: false).

	-o=FILE : string, output file where the generated parser will be
	written (default: stdout).

	-x : boolean, if set, do not build the parser, just parse the input grammar
	(default: false).

	-receiver-name=NAME : string, name of the receiver variable for the generated
	code blocks. Non-initializer code blocks in the grammar end up as methods on the
	*current type, and this option sets the name of the receiver (default: c).

The tool makes no attempt to format the code, nor to detect the
required imports. It is recommended to use goimports to properly generate
the output code:
	pigeon GRAMMAR_FILE | goimports > output_file.go

The goimports tool can be installed with:
	go get golang.org/x/tools/cmd/goimports

If the code blocks in the grammar are golint- and go vet-compliant, then
the resulting generated code will also be golint- and go vet-compliant.

The generated code doesn't use any third-party dependency unless code blocks
in the grammar require such a dependency.

PEG syntax

The accepted syntax for the grammar is formally defined in the
grammar/pigeon.peg file, using the PEG syntax. What follows is an informal
description of this syntax.

Error reporting

*/
package main
