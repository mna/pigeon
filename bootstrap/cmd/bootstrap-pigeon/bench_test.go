package main

import (
	"io/ioutil"
	"testing"
)

func BenchmarkParsePigeonNoMemo(b *testing.B) {
	memoize = false
	d, err := ioutil.ReadFile("../../../grammar/pigeon.peg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParsePigeonMemo(b *testing.B) {
	memoize = true
	d, err := ioutil.ReadFile("../../../grammar/pigeon.peg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d); err != nil {
			b.Fatal(err)
		}
	}
}
