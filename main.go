package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
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

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func codeWalk(rootFilePath string, fileTypes []string, delayTimeInMsChannel chan time.Duration, colorChannel chan bool) {
	var files []string
	var fileContents [][]string
	var currentDelay time.Duration = 200

	err := filepath.Walk(rootFilePath, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		var extension = filepath.Ext(file)
		if contains(fileTypes, extension) {
			content, err := readFile(file)
			if err == nil {
				fileContents = append(fileContents, content)
			}
		}
	}

	color.Set(color.FgHiGreen)

	colors := []color.Attribute{
		color.FgHiGreen,
		color.FgBlue,
		color.FgCyan,
		color.FgHiBlue,
		color.FgHiMagenta,
		color.FgHiWhite,
	}

	rand.Seed(time.Now().Unix())

	defer color.Unset()

	for _, c := range fileContents {
		for _, l := range c {
			for _, x := range l {

				select {
				case delayTimeInMs := <-delayTimeInMsChannel:
					currentDelay = delayTimeInMs
				default:
				}

				select {
				case colorSwitch := <-colorChannel:
					randomColorIndex := rand.Intn(6)
					if colorSwitch {
						color.Set(colors[randomColorIndex])
					}
				default:
				}

				time.Sleep(currentDelay * time.Millisecond)
				fmt.Print(string(x))
			}
			fmt.Println("")
		}
	}
}

func keyHandler(delayChannel chan time.Duration, colorChannel chan bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	var delay time.Duration = 200

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}
			if ev.Ch == '+' {
				delay = delay - 50
				if delay > 0 {
					delayChannel <- delay
				}
			}
			if ev.Ch == '-' {
				delay = delay + 50
				delayChannel <- delay
			}
			if ev.Ch == 'c' {
				colorChannel <- true
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			break mainloop
		}
	}

}

func main() {
	root := "/home/icke/workspace/playground/kubernetes"
	fileTypes := []string{".tf", ".sh", ".java", ".go"}
	delay := make(chan time.Duration)
	color := make(chan bool)
	go codeWalk(root, fileTypes, delay, color)
	keyHandler(delay, color)
}
