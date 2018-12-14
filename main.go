package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func readFile(fileName string){
	file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func main(){

	fileTypes := []string{".tf", ".sh", ".java"}
	var files []string

    root := "/home/icke/workspace/qudo"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
    for _, file := range files {
    	var extension = filepath.Ext(file)
    	if contains(fileTypes, extension){
			fmt.Println(file)
			fmt.Println(extension)
    		readFile(file)
		}
    }
}
