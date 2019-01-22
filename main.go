package main

import (
	"flag"
	"os"
	"path/filepath"
)

const DEFAULT_LOG = "/tmp/code_walk.log"

var homeDir  = os.Getenv("HOME")
var codeWalkDir = homeDir + "/.code_walk"
var fileTypes = []string{".tf", ".sh", ".java", ".go"}


var codeChannel = make(chan string)
var directoryToWalk string

func initCodeWalk() {
	Info.Println("create if not exist: ", codeWalkDir)
	createDirIfNotExist(codeWalkDir)
}

func readSrcFiles(dir string) []string {
	var files []string
	Info.Println("Walking Directory: ", dir)
	err := filepath.Walk(dir, visit(&files))
	if err != nil {
		panic(err)
	}
	return files
}

func main() {
	initLogging(DEFAULT_LOG)

	initCodeWalk()

	dir := flag.String("dir", "/", "directory to walk")
	flag.Parse()
	directoryToWalk = *dir
	files := readSrcFiles(directoryToWalk)


	go engine(files, fileTypes)
	ui()
}
