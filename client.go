package main

import (
	"net/rpc"
)

func client(msg string) {
	var reply string
	client, err := rpc.DialHTTP("tcp", "localhost:8823")
	if err != nil {
		Error.Println("Connection error: ", err)
	} else {
		client.Call("Command.Receive", msg, &reply)
	}
}
