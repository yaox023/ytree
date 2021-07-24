package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const standardText = `testDir
├── a
│   └── 1.txt
└── b
    ├── 2.txt
    └── c
        └── 3.txt
`

const standardJSON = `{
	"type": "directory",
	"name": "testDir",
	"contents": [
		{
			"type": "directory",
			"name": "a",
			"contents": [
				{
					"type": "file",
					"name": "1.txt"
				}
			]
		},
		{
			"type": "directory",
			"name": "b",
			"contents": [
				{
					"type": "file",
					"name": "2.txt"
				},
				{
					"type": "directory",
					"name": "c",
					"contents": [
						{
							"type": "file",
							"name": "3.txt"
						}
					]
				}
			]
		}
	]
}`

func init() {
	os.MkdirAll("./testResult", 0755)
}

func TestJSON(t *testing.T) {
	outputPath := "./testResult/test_json_output.json"
	// templatePath := "./testResult/test_json_template.json"

	os.Args = append(os.Args, "-dir=./testDir")
	os.Args = append(os.Args, "-json")
	os.Args = append(os.Args, fmt.Sprintf("-output=%s", outputPath))
	main()
	os.Args = os.Args[:len(os.Args)-3]

	outputBytes, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Error(err)
	}

	if string(outputBytes) != standardJSON {
		t.Error("result not match")
	}
}

func TestText(t *testing.T) {
	outputPath := "./testResult/test_text_output.txt"
	// templatePath := "./testResult/test_text_template.txt"

	os.Args = append(os.Args, "-dir=./testDir")
	os.Args = append(os.Args, fmt.Sprintf("-output=%s", outputPath))
	main()
	os.Args = os.Args[:len(os.Args)-2]

	outputBytes, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Error(err)
	}

	if string(outputBytes) != standardText {
		t.Error("result not match")
	}
}

func TestStdout(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = append(os.Args, "-dir=./testDir")
	main()
	os.Args = os.Args[:len(os.Args)-1]

	w.Close()
	os.Stdout = rescueStdout
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	// stdout output a extra \n
	if string(out) != standardText+"\n" {
		t.Error("result not match")
	}
}
