package main

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	os.Args = append(os.Args, "-dir=./testDir")
	main()
}
