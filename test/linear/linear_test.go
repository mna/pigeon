package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestLinearTime(t *testing.T) {
	var buf bytes.Buffer

	if testing.Short() {
		t.Skip()
	}

	sizes := []int64{
		1 << 10,   // 1Kb
		10 << 10,  // 10Kb
		100 << 10, // 100Kb
		1 << 20,   // 1MB
	}
	for _, sz := range sizes {
		buf.Reset()
		r := io.LimitReader(rand.Reader, sz)
		enc := base64.NewEncoder(base64.StdEncoding, &buf)
		_, err := io.Copy(enc, r)
		if err != nil {
			t.Fatal(err)
		}
		enc.Close()

		if testing.Verbose() {
			fmt.Printf("starting with %dKB...\n", sz/1024)
		}
		start := time.Now()
		if _, err := Parse("", buf.Bytes(), Memoize(true)); err != nil {
			t.Fatal(err)
		}
		t.Logf("%dKB: %s", sz/1024, time.Now().Sub(start))
	}
}
