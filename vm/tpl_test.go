package vm

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/pigeon/bootstrap"
)

func TestTemplate(t *testing.T) {
	gr, err := bootstrap.NewParser().Parse("", strings.NewReader(`
	{package abc}
	A = !{return true, nil} l1:. l2:'B'i &{return false, nil} l3:[^A-CE\p{Latin}]i {return "ok", nil}
	`))
	if err != nil {
		t.Fatal(err)
	}

	pg, err := NewGenerator(ioutil.Discard).toProgram(gr)
	if err != nil {
		t.Fatal(err)
	}

	out := ioutil.Discard
	if testing.Verbose() {
		out = os.Stdout
	}
	if err := tpl.Execute(out, pg); err != nil {
		t.Fatal(err)
	}
}
