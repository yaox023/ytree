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

// func outputItemToXML(item *Item, outputPath string) {

// }

func writeFile(data []byte, filePath string) error {
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	dir := flags.String("dir", "./testDir", "directory to show")
	toJSON := flags.Bool("json", false, "output to a json file")
	toXML := flags.Bool("xml", false, "output to a xml file")
	outputPath := flags.String("output", "", "output path")
	flags.Parse(os.Args[2:])

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
		jsonBytes, err := outputItemToJSON(item, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
		err = writeFile(jsonBytes, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *toXML {
		// outputItemToXML(item, *outputPath)
		return
	}

	textbytes := outputItemToText(item)
	if *outputPath != "" {
		err = writeFile(textbytes, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Println(string(textbytes))
}
