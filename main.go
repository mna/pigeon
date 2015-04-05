package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/exp/peg/ast"
	"github.com/PuerkitoBio/exp/peg/builder"
)

func main() {
	// define command-line flags
	var (
		dbgFlag        = flag.Bool("debug", false, "set debug mode")
		noBuildFlag    = flag.Bool("x", false, "do not build, only parse")
		outputFlag     = flag.String("o", "", "output file, defaults to stdout")
		curRecvrNmFlag = flag.String("receiver-name", "c", "receiver name for the `current` type's generated methods")
	)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	// get input source
	infile := ""
	if flag.NArg() == 1 {
		infile = flag.Arg(0)
	}
	nm, rc := input(infile)
	defer rc.Close()

	// parse input
	debug = *dbgFlag
	g, err := Parse(nm, rc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error: ", err)
		os.Exit(3)
	}

	if !*noBuildFlag {
		out := output(*outputFlag)
		defer out.Close()

		curNmOpt := builder.CurrentReceiverName(*curRecvrNmFlag)
		if err := builder.BuildParser(out, g.(*ast.Grammar), curNmOpt); err != nil {
			fmt.Fprintln(os.Stderr, "build error: ", err)
			os.Exit(5)
		}
	}
}

func usage() {
	fmt.Printf("usage: %s [options] FILE\n", os.Args[0])
	flag.PrintDefaults()
}

func input(filename string) (nm string, rc io.ReadCloser) {
	nm = "stdin"
	inf := os.Stdin
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		inf = f
		nm = filename
	}
	r := bufio.NewReader(inf)
	return nm, makeReadCloser(r, inf)
}

func output(filename string) io.WriteCloser {
	out := os.Stdout
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
		out = f
	}
	return out
}

func makeReadCloser(r io.Reader, c io.Closer) io.ReadCloser {
	rc := struct {
		io.Reader
		io.Closer
	}{r, c}
	return io.ReadCloser(rc)
}

func (c *current) astPos() ast.Pos {
	return ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
