package main

import (
	"net/rpc"
)

func sendTo(info CodeWalkFileInfo) {
	var reply string
	client, err := rpc.DialHTTP("tcp", "localhost:8823")
	if err != nil {
		Error.Println("Connection error: ", err)
	} else {
		callErr := client.Call("Command.Receive", info, &reply)
		if callErr != nil {
			Error.Println(callErr)
		}
	}
}
