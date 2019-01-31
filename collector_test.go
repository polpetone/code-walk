package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetCommitDates(t *testing.T) {
	//TODO: get resilient path
	///playground/go/src/crypto/rand/rand.go
	first, last, err :=  getCommitDates("/home/icke/workspace/playground/go/src/crypto/rand/rand.go")
	fmt.Println(err)
	layout := "2006-01-02 15:04:05 +0100"

	firstTime, _ := time.Parse(layout, first.String())
	lastTime, err := time.Parse(layout, last.String())

	fmt.Println(err)

	fmt.Println("first", first)
	fmt.Println("last", last)
	fmt.Println("first", firstTime)
	fmt.Println("last", lastTime)

	distance := last.Sub(first)
	distance1 := first.Sub(last)

	fmt.Println("distance ", distance)
	fmt.Println("distance1 ", distance1)
}

