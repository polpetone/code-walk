package main

import (
	"os"
)

const DEFAULT_LOG = "/tmp/code_walk.log"

var homeDir  = os.Getenv("HOME")
var codeWalkDir = homeDir + "/.code_walk"


func initCodeWalk() {
	Info.Println("create if not exist: ", codeWalkDir)
	createDirIfNotExist(codeWalkDir)
}

func main() {
	initLogging(DEFAULT_LOG)
	initCodeWalk()

	engine()
}
