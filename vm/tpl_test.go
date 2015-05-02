package vm

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/pigeon/bootstrap"
)

func TestTemplate(t *testing.T) {
	gr, err := bootstrap.NewParser().Parse("", strings.NewReader("A = 'a'"))
	if err != nil {
		t.Fatal(err)
	}

	pg, err := NewGenerator(ioutil.Discard).toProgram(gr)
	if err != nil {
		t.Fatal(err)
	}

	if err := tpl.Execute(os.Stdout, pg); err != nil {
		t.Fatal(err)
	}
}
