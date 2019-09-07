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
		pgot, err := ParseFile(file)
		if err != nil {
			t.Errorf("%s: pigeon.ParseFile: %v", file, err)
			continue
		}

		pogot, err := optimized.ParseFile(file)
		if err != nil {
			t.Errorf("%s: optimized.ParseFile: %v", file, err)
			continue
		}

		poggot, err := optimizedgrammar.ParseFile(file)
		if err != nil {
			t.Errorf("%s: optimizedgrammar.ParseFile: %v", file, err)
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

		if !reflect.DeepEqual(pogot, jgot) {
			t.Errorf("%s: optimized not equal", file)
			continue
		}

		if !reflect.DeepEqual(poggot, jgot) {
			t.Errorf("%s: optimized grammar not equal", file)
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

func TestChoiceAltStatistics(t *testing.T) {
	cases := []struct {
		json          string
		expectedStats map[string]map[string]int
	}{
		{
			json: `{}`,
			expectedStats: map[string]map[string]int{
				"Value 21:15": {
					"1": 1,
				},
			},
		},
		{
			json: `{ "string": "string", "number": 123 }`,
			expectedStats: map[string]map[string]int{
				"Integer 60:11": {
					"2":        1,
					"no match": 1,
				},
				"String 64:16": {
					"1":        18,
					"no match": 3,
				},
				"Value 21:15": {
					"1": 1,
					"3": 1,
					"4": 1,
				},
			},
		},
	}

	for _, test := range cases {
		stats := Stats{}
		_, err := Parse("TestStatistics", []byte(test.json), Statistics(&stats, "no match"))
		if err != nil {
			t.Fatalf("Expected to parse %s without error, got: %v", test.json, err)
		}
		if !reflect.DeepEqual(test.expectedStats, stats.ChoiceAltCnt) {
			t.Fatalf("Expected stats to equal %#v, got %#v", test.expectedStats, stats.ChoiceAltCnt)
		}
	}
}

func TestZeroZero(t *testing.T) {
	_, err := Parse(`fuzz`, []byte(`00`))
	if err == nil {
		t.Errorf("Error expected for invalid input: `00`")
	}
}

// The forward slash (solidus) is not a valid escape in Go, it will normally
// fail if there's one in the string
func TestForwardSlash(t *testing.T) {
	_, err := Parse(`fuzz`, []byte(`"\/"`))
	if err != nil {
		t.Errorf("forward slash pigeon.Parse: %v", err)
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
