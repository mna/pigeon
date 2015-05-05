package main

import (
	"io/ioutil"
	"log"

	"github.com/PuerkitoBio/pigeon/vm"
	"github.com/davecheney/profile"
)

func main() {
	d, err := ioutil.ReadFile("../../../grammar/pigeon.peg")
	if err != nil {
		log.Fatal(err)
	}
	defer profile.Start(profile.CPUProfile).Stop()

	// TODO : that doesn't work, no program to run...
	if _, err := vm.Parse("", d, vm.Memoize(false)); err != nil {
		log.Fatal(err)
	}
}
