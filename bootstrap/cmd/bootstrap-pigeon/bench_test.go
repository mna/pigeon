package main

import (
	"io/ioutil"
	"testing"
)

func BenchmarkParsePigeon(b *testing.B) {
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
