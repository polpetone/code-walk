package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func readFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return text, nil
}

func appendLineToFile(fileName string, content string) error {
	var file, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	if err == nil {
		defer file.Close()
		_, err = file.WriteString(content)
		_, err = file.WriteString("\n")
		if err == nil {
			err = file.Sync()
		}
	}

	return err
}

func createFile(filePath string) error {
	var _, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if err == nil {
			defer file.Close()
			return nil
		}
	}
	return err
}


func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func cutFileNameFromPath(path string) (string, string){

	i := strings.LastIndex(path, "/")
	fileName := path[i+1:]
	dir := path[:i]

	return fileName, dir
}

func removeWordFromString(s string, word string) string {
	return strings.TrimSpace(strings.TrimPrefix(s, word))
}

func keysFromMap(m map[string][]string) (keys []string){
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}


func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}
