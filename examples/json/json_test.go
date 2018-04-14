package json

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	optimized "github.com/mna/pigeon/examples/json/optimized"
	optimizedgrammar "github.com/mna/pigeon/examples/json/optimized-grammar"
)

func TestCmpStdlib(t *testing.T) {
	files := testJSONFiles(t)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			pgot, err := ParseFile(file)
			if err != nil {
				t.Errorf("%s: pigeon.ParseFile: %v", file, err)
				return
			}

			pogot, err := optimized.ParseFile(file)
			if err != nil {
				t.Errorf("%s: optimized.ParseFile: %v", file, err)
				return
			}

			poggot, err := optimizedgrammar.ParseFile(file)
			if err != nil {
				t.Errorf("%s: optimizedgrammar.ParseFile: %v", file, err)
				return
			}

			b, err := ioutil.ReadFile(file)
			if err != nil {
				t.Errorf("%s: ioutil.ReadAll: %v", file, err)
				return
			}
			var jgot interface{}
			if err := json.Unmarshal(b, &jgot); err != nil {
				t.Errorf("%s: json.Unmarshal: %v", file, err)
				return
			}

			if !reflect.DeepEqual(pgot, jgot) {
				t.Errorf("%s: not equal", file)
				return
			}

			if !reflect.DeepEqual(pogot, jgot) {
				t.Errorf("%s: optimized not equal", file)
				return
			}

			if !reflect.DeepEqual(poggot, jgot) {
				t.Errorf("%s: optimized grammar not equal", file)
				return
			}
		})
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
		t.Run(fi.Name(), func(t *testing.T) {
			if filepath.Ext(fi.Name()) == ".json" {
				files = append(files, filepath.Join(rootDir, fi.Name()))
			}
		})
	}
	return files
}

func TestChoiceAltStatistics(t *testing.T) {
	cases := []struct {
		name          string
		json          string
		expectedStats map[string]map[string]int
	}{
		{
			name: "empty json",
			json: `{}`,
			expectedStats: map[string]map[string]int{
				"Bool 92:8": {
					"no match": 1,
				},
				"Integer 68:11": {
					"no match": 1,
				},
				"Value 29:15": {
					"1":        1,
					"no match": 1,
				},
			},
		},
		{
			name: "simple json",
			json: `{ "string": "string", "number": 123 }`,
			expectedStats: map[string]map[string]int{
				"Integer 68:11": {
					"2":        1,
					"no match": 2,
				},
				"Bool 92:8": {
					"no match": 1,
				},
				"String 72:16": {
					"1":        18,
					"no match": 3,
				},
				"Value 29:15": {
					"1":        1,
					"3":        1,
					"4":        1,
					"no match": 1,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			stats := Stats{}
			_, err := Parse("TestStatistics", []byte(test.json), Statistics(&stats, "no match"))
			if err != nil {
				t.Fatalf("Expected to parse %s without error, got: %v", test.json, err)
			}
			if !reflect.DeepEqual(test.expectedStats, stats.ChoiceAltCnt) {
				t.Fatalf("Expected stats to equal %#v, got %#v", test.expectedStats, stats.ChoiceAltCnt)
			}
		})
	}
}

func BenchmarkPigeonJSONNoMemo(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d, Memoize(false)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPigeonJSONMemo(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d, Memoize(true)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPigeonJSONOptimized(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := optimized.Parse("", d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPigeonJSONOptimizedGrammar(b *testing.B) {
	d, err := ioutil.ReadFile("testdata/github-octokit-repos.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := optimizedgrammar.Parse("", d); err != nil {
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
