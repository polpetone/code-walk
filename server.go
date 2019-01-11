package main

import (
	"net"
	"net/http"
	"net/rpc"
)

func server(){
	command := new(Command)
	err := rpc.Register(command)
	if err != nil {
		Error.Println("Format of service Command isn't correct", err)
	}
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":8123")
	if e != nil {
		Error.Println("Listen error: ", e)
	}
	Info.Println("Serving RPC server on port ", 8123)
	err = http.Serve(listener, nil)
	if err != nil {
		Error.Println("Error serving: ", err)
	}
}
