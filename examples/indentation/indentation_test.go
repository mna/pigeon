package main

import (
	"testing"
    "io/ioutil"
	"path/filepath"
)

func TestIndentation(t *testing.T) {
	files := testIndentationFiles(t)
	for _, file := range files {
		pgot, err := ParseFile(file)
		if err != nil {
			t.Errorf("%s: pigeon.ParseFile: %v", file, err)
			continue
		}
        got,err := pgot.(ProgramNode).exec()
        if err != nil {
			t.Errorf("%s: ProgramNode.exec: %v", file, err)
			continue
		}
        
        exp := 42
		if got != exp {
			t.Errorf("%v: want %v, got %v", file, exp, got)
		}
	}
}

func testIndentationFiles(t *testing.T) []string {
	const rootDir = "testdata"

	fis, err := ioutil.ReadDir(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 0, len(fis))
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) == ".txt" {
			files = append(files, filepath.Join(rootDir, fi.Name()))
		}
	}
	return files
}
