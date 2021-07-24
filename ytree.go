package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
)

type ItemType string

const (
	fileItem ItemType = "file"
	dirItem  ItemType = "directory"
)

type Item struct {
	XMLName  xml.Name `json:"-"`
	Type     ItemType `json:"type" xml:"-"`
	Name     string   `json:"name" xml:"name,attr"`
	Contents []Item   `json:"contents,omitempty"`
}

type Report struct {
	XMLName     xml.Name `xml:"report"`
	Directories int      `xml:"directories"`
	Files       int      `xml:"files"`
}

func readItem(itemPath string, item *Item, report *Report) {
	itemInfo, err := os.Stat(itemPath)
	if err != nil {
		log.Fatal(err)
	}
	if !itemInfo.IsDir() {
		item.Type = fileItem
		item.Name = itemInfo.Name()
		item.XMLName = xml.Name{Local: string(fileItem)}
		report.Files += 1
		return
	}
	fileInfos, err := ioutil.ReadDir(itemPath)
	if err != nil {
		log.Fatal(err)
	}
	item.Type = dirItem
	item.Name = itemInfo.Name()
	item.XMLName = xml.Name{Local: string(dirItem)}
	item.Contents = []Item{}
	report.Directories += 1
	for _, fileInfo := range fileInfos {
		newItem := &Item{}
		readItem(filepath.Join(itemPath, fileInfo.Name()), newItem, report)
		item.Contents = append(item.Contents, *newItem)
	}
}

func printItem(item *Item, buffer *bytes.Buffer, levels []bool) {

	for index, level := range levels {
		if index == len(levels)-1 {
			if level {
				buffer.WriteString("└── ")
			} else {
				buffer.WriteString("├── ")
			}
		} else {
			if level {
				buffer.WriteString("    ")
			} else {
				buffer.WriteString("│   ")
			}
		}
	}

	buffer.WriteString(item.Name)
	buffer.WriteString("\n")

	length := len(item.Contents)
	for index, childItem := range item.Contents {
		if index == length-1 {
			printItem(&childItem, buffer, append(levels, true))
		} else {
			printItem(&childItem, buffer, append(levels, false))
		}
	}
}

func outputItemToText(item *Item) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(item.Name)
	buffer.WriteString("\n")

	length := len(item.Contents)
	for index, childItem := range item.Contents {
		if index == length-1 {
			printItem(&childItem, &buffer, []bool{true})
		} else {
			printItem(&childItem, &buffer, []bool{false})
		}
	}
	return buffer.Bytes()
}

func outputItemToJSON(item *Item, outputPath string) ([]byte, error) {
	output, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return output, nil
}

func outputItemToXML(item *Item, report *Report, outputPath string) ([]byte, error) {
	tree := struct {
		XMLName xml.Name
		Item    Item
		Report  Report
	}{
		xml.Name{Local: "tree"},
		*item,
		*report,
	}
	output, err := xml.MarshalIndent(tree, "", "  ")
	if err != nil {
		return nil, err
	}
	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	output = []byte(header + string(output))
	return output, nil
}

func writeFile(data []byte, filePath string) error {
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	parser := argparse.NewParser("ytree", "list contents of directories in a tree-like format")
	dir := parser.String("d", "dir", &argparse.Options{Default: ".", Help: "directory to list"})
	toJSON := parser.Flag("J", "json", &argparse.Options{Help: "Turn on JSON output. Outputs the directory tree as an JSON formatted array."})
	toXML := parser.Flag("X", "xml", &argparse.Options{Help: "Turn on XML output. Outputs the directory tree as an XML formatted file."})
	outputPath := parser.String("o", "filename", &argparse.Options{Help: "Send output to filename."})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	dirInfo, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("Input is not directory: ", *dir)
	}

	item := &Item{}
	// -1 for subtrack start dir
	report := &Report{Files: 0, Directories: -1}
	readItem(*dir, item, report)

	var outputBytes []byte

	if *toJSON {
		outputBytes, err = outputItemToJSON(item, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
	} else if *toXML {
		outputBytes, err = outputItemToXML(item, report, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		outputBytes = outputItemToText(item)
	}

	if *outputPath != "" {
		err = writeFile(outputBytes, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(string(outputBytes))
	}
}
