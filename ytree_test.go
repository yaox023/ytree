package main

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	os.Args = append(os.Args, "-dir=./")
	// os.Args = append(os.Args, "-json")
	// os.Args = append(os.Args, "-output=output.json")
	main()
}

/*
├── README.md
├── go.mod
├── output.json
├── testDir
│   ├── a
│       └── 1.txt
│   └── b
│       ├── 2.txt
│       └── c
│           └── 3.txt
├── ytree.go
└── ytree_test.go

├── README.md
├── go.mod
├── output.json
├── testDir
│   ├── a
│   │   └── 1.txt
│   └── b
│       ├── 2.txt
│       └── c
│           └── 3.txt
├── ytree.go
└── ytree_test.go
*/
