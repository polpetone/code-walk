package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

const DEFAULT_LOG = "/tmp/code_walk.log"

var homeDir  = os.Getenv("HOME")
var codeWalkDir = homeDir + "/.code_walk"
var fileTypes = []string{".tf", ".sh", ".java", ".go"}


type Command int


func (c *Command) Receive(key string, reply *string) error {
	fmt.Println("received key", key)
	return nil
}


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

	command := new(Command)

	err := rpc.Register(command)

	if err != nil {
		fmt.Println("Format of service Command isn't correct", err)
	}
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":8123")
	if e != nil {
		fmt.Println("Listen error: ", e)
	}
	fmt.Println("Serving RPC server on port ", 8123)
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println("Error serving: ", err)
	}

	//dir := flag.String("dir", "/", "directory to walk")
	//flag.Parse()
	//files := readSrcFiles(*dir)

	//engine(files, fileTypes)
}
