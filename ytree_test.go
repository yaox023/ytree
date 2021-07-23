package main

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	os.Args = append(os.Args, "-dir=./testDir")
	// os.Args = append(os.Args, "-json")
	// os.Args = append(os.Args, "-output=output.json")
	main()
}
