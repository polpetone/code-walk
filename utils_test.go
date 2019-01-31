package main

import (
	"fmt"
	"testing"
)

func TestCutFilenameFromPath(t *testing.T) {

	path := "/foo/bar/file.go"
	fileName, dir := cutFileNameFromPath(path)

	//TODO: asserts + more test cases
	fmt.Println(fileName)
	fmt.Println(dir)
}

func TestRemoveWordFromString(t *testing.T){
	s := removeWordFromString("author foo bar", "author")
	fmt.Println(s)
}
