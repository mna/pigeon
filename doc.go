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

Identifiers, whitespace, comments and literals follow the same
notation as the Go language, as defined in the language specification
(http://golang.org/ref/spec#Source_code_representation):

	// single line comment*/
//	/* multi-line comment */
/*	'x' (single quotes for single char literal)
	"double quotes for string literal"
	`backtick quotes for raw string literal`
	RuleName (a valid identifier)

The grammar must be Unicode text encoded in UTF-8. New lines are identified
by the \n character (U+000A). Space (U+0020), horizontal tabs (U+0009) and
carriage returns (U+000D) are considered whitespace and are ignored except
to separate tokens.

Rules

A PEG grammar consists of a set of rules. A rule is an identifier followed
by a rule definition operator and an expression. An optional display name -
a string literal used in error messages instead of the rule identifier - can
be specified after the rule identifier. E.g.:
	RuleA "friendly name" = 'a'+ // RuleA is one or more lowercase 'a's

The rule definition operator can be any one of those:
	=, <-, ← (U+2190), ⟵ (U+27F5)

Expressions

A rule is defined by an expression. The following sections describe the
various expression types. Expressions can be grouped by using parentheses,
and a rule can be referenced by its identifier in place of an expression.

Choice expression

The choice expression is a list of expressions that will be tested in the
order they are defined. The first one that matches will be used. Expressions
are separated by the forward slash character "/". E.g.:
	ChoiceExpr = A / B / C // A, B and C should be rules declared in the grammar

Because the first match is used, it is important to think about the order
of expressions. For example, in this rule, "<=" would never be used because
the "<" expression comes first:
	BadChoiceExpr = "<" / "<="

Sequence expression

The sequence expression is a list of expressions that must all match in
that same order for the sequence expression to be considered a match.
Expressions are separated by whitespace. E.g.:
	SeqExpr = "A" "b" "c" // matches "Abc", but not "Acb"

Labeled expression

A labeled expression consists of an identifier followed by a colon ":"
and an expression. A labeled expression introduces a variable named with
the label that can be referenced in the parent expression's code block.
The variable will have the value of the expression that follows the colon.
E.g.:
	LabeledExpr = value:[a-z]+ {
		fmt.Println(value)
		return value, nil
	}

And and not expressions

An expression prefixed with the ampersand "&" is the "and" predicate
expression: it is considered a match if the following expression is a match,
but it does not consume any input.

An expression prefixed with the exclamation point "!" is a predicate
expression: it is considered a match if the following expression is not
a match, but it does not consume any input. E.g.:
	AndExpr = "A" &"B" // matches "A" if followed by a "B" (does not consume "B")
	NotExpr = "A" !"B" // matches "A" if not followed by a "B" (does not consume "B")

The expression following the & and ! operators can be a code block. In that
case, the code block must return a bool and an error. The operator's semantic
is the same, & is a match if the code block returns true, ! is a match if the
code block returns false. The code block has access to any labeled value
defined in its scope. E.g.:
	CodeAndExpr = value:[a-z] &{
		// can access the value local variable...
		return true, nil
	}

Repeating expressions

An expression followed by "*", "?" or "+" is a match if the expression
occurs zero or more times ("*"), zero or one time "?" or one or more times
("+") respectively. The match is greedy, it will match as many times as
possible. E.g.
	ZeroOrMoreAs = "A"*

Literal matcher

A literal matcher tries to match the input against a single character or a
string literal. The literal may be a single-quoted single character, a
double-quoted string or a backtick-quoted raw string. The same rules as in Go
apply regarding the allowed characters and escapes.

The literal may be followed by a lowercase "i" (outside the ending quote)
to indicate that the match is case-insensitive. E.g.:
	LiteralMatch = "Awesome\n"i // matches "awesome" followed by a newline

Character class matcher

A character class matcher tries to match the input against a class of characters
inside square brackets "[...]". Inside the brackets, characters represent
themselves and the same escapes as in string literals are available, except
that the single- and double-quote escape is not valid, instead the closing
square bracket "]" must be escaped to be used.

Character ranges can be specified using the "[a-z]" notation. Unicode
classes can be specified using the "[\pL]" notation, where L is a
single-letter Unicode class of characters, or using the "[\p{Class}]"
notation where Class is a valid Unicode class (e.g. "Latin").

As for string literals, a lowercase "i" may follow the matcher (outside
the ending square bracket) to indicate that the match is case-insensitive.
A "^" as first character inside the square brackets indicates that the match
is inverted (it is a match if the input does not match the character class
matcher). E.g.:
	NotAZ = [^a-z]i

Any matcher

The any matcher is represented by the dot ".". It matches any character
except the end of file, thus the "!." expression is used to indicate "match
the end of file". E.g.:
	AnyChar = . // match a single character
	EOF = !.

Code block

Code blocks can be added to generate custom Go code. There are three kinds
of code blocks: the initializer, the action and the predicate. All code blocks
appear inside curly braces "{...}".

The initializer must appear first in the grammar, before any rule. It is
copied as-is (minus the wrapping curly braces) at the top of the generated
parser. It may contain function declarations, types, variables, etc. just
like any Go file. Every symbol declared here will be available to all other
code blocks.  Although the initializer is optional in a valid grammar, it is
usually required to generate a valid Go source code file (for the package
clause). E.g.:
	{
		package main

		func someHelper() {
			// ...
		}
	}

Action code blocks are code blocks declared after an expression in a rule.
Those code blocks are turned into a method on the "*current" type in the
generated source code. The method receives any labeled expression's value
as argument (as interface{}) and must return two values, the first being
the value of the expression (an interface{}), and the second an error.
If a non-nil error is returned, it is added to the list of errors that the
parser will return. E.g.:
	RuleA = "A"+ {
		// return the matched string, "c" is the default name for
		// the *current receiver variable.
		return string(c.text), nil
	}

Predicate code blocks are code blocks declared immediately after the and "&"
or the not "!" operators. Like action code blocks, predicate code blocks
are turned into a method on the "*current" type in the generated source code.
The method receives any labeled expression's value as argument (as interface{})
and must return two values, the first being a bool and the second an error.
If a non-nil error is returned, it is added to the list of errors that the
parser will return. E.g.:
	RuleAB = [ab]i+ &{
		return true, nil
	}

The current type is a struct that provides two useful fields that can be
accessed in action and predicate code blocks: "pos" and "text".

The "pos" field indicates the current position of the parser in the source
input. It is itself a struct with three fields: "line", "col" and "offset".
Line is a 1-based line number, col is a 1-based column number that counts
runes from the start of the line, and offset is a 0-based byte offset.

The "text" field is the slice of bytes of the current match. It is empty
in a predicate code block.

Using the generated parser

The parser generated by pigeon exports a few symbols so that it can be used
as a package with public functions to parse input text.

TODO: Start rule is the first rule. Example package to document the exported symbols.

Error reporting

TODO: List of errors, ParserError type, grammar example to handle common error (like
non-terminated string literal), panic.

*/
package main
