package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCmpStdlib(t *testing.T) {
	files := testJSONFiles(t)
	for _, file := range files {
		pgot, err := ParseFile(file)
		if err != nil {
			t.Errorf("%s: pigeon.ParseFile: %v", file, err)
			continue
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("%s: ioutil.ReadAll: %v", file, err)
			continue
		}
		var jgot interface{}
		if err := json.Unmarshal(b, &jgot); err != nil {
			t.Errorf("%s: json.Unmarshal: %v", file, err)
			continue
		}

		if !reflect.DeepEqual(pgot, jgot) {
			t.Errorf("%s: not equal", file)
			continue
		}
	}
}

func testJSONFiles(t *testing.T) []string {
	const rootDir = "testdata"

	fis, err := ioutil.ReadDir(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 0, len(fis))
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) == ".json" {
			files = append(files, filepath.Join(rootDir, fi.Name()))
		}
	}
	return files
}

func BenchmarkPigeonJSON(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	br := bytes.NewReader(d)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", br); err != nil {
			b.Fatal(err)
		}
		if _, err := br.Seek(0, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdlibJSON(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var iface interface{}
		if err := json.Unmarshal(d, &iface); err != nil {
			b.Fatal(err)
		}
	}
}
