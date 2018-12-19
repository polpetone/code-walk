package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"github.com/fatih/color"
)


func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func readFile(fileName string) ([]string, error){
	file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    var text []string
    for scanner.Scan() {
        text = append(text,scanner.Text())
    }
    if err := scanner.Err(); err != nil {
    	return nil, err
    }
    return text, nil
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func printContent(contents [][]string, delayTimeInMS time.Duration){
	color.Set(color.FgHiGreen)
	defer color.Unset()
	for _,c := range contents {
		for _, l := range c {
			time.Sleep(delayTimeInMS * time.Millisecond)
			fmt.Println(l)
		}
	}
}

func main(){

	fileTypes := []string{".tf", ".sh", ".java"}
	var files []string
	var fileContents [][]string


    root := "/home/icke/workspace/qudo"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
    for _, file := range files {
    	var extension = filepath.Ext(file)
    	if contains(fileTypes, extension){
    		content, err := readFile(file)
    		if err == nil {
    			fileContents = append(fileContents, content)
			}
		}
    }
    printContent(fileContents, 200)
}
