package ytree

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	os.Args = append(os.Args, "-dir=./testDir")

	print()
}
