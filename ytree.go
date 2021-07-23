package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ItemType string

const (
	fileItem ItemType = "file"
	dirItem  ItemType = "directory"
)

type Item struct {
	Type     ItemType `json:"type"`
	Name     string   `json:"name"`
	Contents []Item   `json:"contents,omitempty"`
}

func readItem(itemPath string, item *Item) {
	itemInfo, err := os.Stat(itemPath)
	if err != nil {
		log.Fatal(err)
	}
	if !itemInfo.IsDir() {
		item.Type = fileItem
		item.Name = itemInfo.Name()
		return
	}
	fileInfos, err := ioutil.ReadDir(itemPath)
	if err != nil {
		log.Fatal(err)
	}
	item.Type = dirItem
	item.Name = itemInfo.Name()
	item.Contents = []Item{}
	for _, fileInfo := range fileInfos {
		newItem := &Item{}
		readItem(filepath.Join(itemPath, fileInfo.Name()), newItem)
		item.Contents = append(item.Contents, *newItem)
	}
}

/*
├── README.md
├── go.mod
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

// TODO still have some bug, for parent |
func printItem(item *Item, buffer *bytes.Buffer, level int, isLast bool, isRootLast bool) {
	originalLevel := level
	for level > 0 {
		if level == originalLevel {
			if isRootLast {
				buffer.WriteString("    ")
			} else {
				buffer.WriteString("│   ")
			}
		} else {
			buffer.WriteString("    ")
		}
		level -= 1
	}

	if isLast {
		buffer.WriteString("└── ")
	} else {
		buffer.WriteString("├── ")
	}

	buffer.WriteString(item.Name)
	buffer.WriteString("\n")

	length := len(item.Contents)
	for index, childItem := range item.Contents {
		if index == length-1 {
			printItem(&childItem, buffer, originalLevel+1, true, isRootLast)
		} else {
			printItem(&childItem, buffer, originalLevel+1, false, isRootLast)
		}
	}
}

func outputItemToConsole(item *Item) {
	var buffer bytes.Buffer
	buffer.WriteString(item.Name)
	buffer.WriteString("\n")

	length := len(item.Contents)
	for index, childItem := range item.Contents {
		if index == length-1 {
			printItem(&childItem, &buffer, 0, true, true)
		} else {
			printItem(&childItem, &buffer, 0, false, false)
		}
	}
	fmt.Println(buffer.String())
}

func outputItemToJSON(item *Item, outputPath string) {
	output, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(outputPath, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func outputItemToXML(item *Item, outputPath string) {

}

func main() {
	dir := flag.String("dir", "./", "directory to show")
	toJSON := flag.Bool("json", false, "output to a json file")
	toXML := flag.Bool("xml", false, "output to a xml file")
	outputPath := flag.String("output", "", "output path")
	flag.Parse()

	dirInfo, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("Input is not directory: ", *dir)
	}

	item := &Item{}
	readItem(*dir, item)

	if *toJSON {
		outputItemToJSON(item, *outputPath)
		return
	}

	if *toXML {
		outputItemToXML(item, *outputPath)
		return
	}

	outputItemToConsole(item)

	// output to txt file

	// readDir(*dir)
	// fmt.Println(*dir)
	// printDir(*dir, 2)
}
