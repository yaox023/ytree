package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

const standardXML = `<?xml version="1.0" encoding="UTF-8"?>
<tree>
  <directory name="testDir">
    <directory name="a">
      <file name="1.txt"></file>
    </directory>
    <directory name="b">
      <file name="2.txt"></file>
      <directory name="c">
        <file name="3.txt"></file>
      </directory>
    </directory>
  </directory>
  <report>
    <directories>3</directories>
    <files>3</files>
  </report>
</tree>`

func init() {
	os.MkdirAll("./testResult", 0755)

	// delete test args
	var newArgs []string
	for _, arg := range os.Args {
		if !strings.HasPrefix(arg, "-test.") {
			newArgs = append(newArgs, arg)
		}
	}
	os.Args = newArgs
}

func TestXML(t *testing.T) {
	outputPath := "./testResult/test_xml_output.xml"

	os.Args = append(os.Args, "-d=./testDir")
	os.Args = append(os.Args, "-X")
	os.Args = append(os.Args, fmt.Sprintf("-o=%s", outputPath))
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

	if string(outputBytes) != standardXML {
		t.Error("result not match")
	}
}

func TestJSON(t *testing.T) {
	outputPath := "./testResult/test_json_output.json"

	os.Args = append(os.Args, "-d=./testDir")
	os.Args = append(os.Args, "-J")
	os.Args = append(os.Args, fmt.Sprintf("-o=%s", outputPath))
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

	os.Args = append(os.Args, "-d=./testDir")
	os.Args = append(os.Args, fmt.Sprintf("-o=%s", outputPath))
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

	os.Args = append(os.Args, "-d=./testDir")
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
