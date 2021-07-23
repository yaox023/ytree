package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var allDirInfo = make(map[string][]os.FileInfo)

func readDir(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	allDirInfo[dir] = append(allDirInfo[dir], fileInfos...)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			readDir(filepath.Join(dir, fileInfo.Name()))
		}
	}
}

func printTab(tab int) {
	fmt.Print("├")
	for i := 0; i < tab; i++ {
		fmt.Print("─")
	}
}

func printDir(base string, tab int) {
	for _, info := range allDirInfo[base] {
		printTab(tab)
		fmt.Println(info.Name())
		if info.IsDir() {
			path := filepath.Join(base, info.Name())
			printDir(path, tab+2)
		}
	}
}

func main() {
	dir := flag.String("dir", "./", "directory to show")
	flag.Parse()

	dirInfo, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("Input is not directory: ", *dir)
	}

	readDir(*dir)
	fmt.Println(*dir)
	printDir(*dir, 2)
}
