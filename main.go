package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"math/rand"
	"path/filepath"
	"time"
)

const DEFAULT_LOG = "/tmp/code_walk.log"
const DELAY_STEP = 20000

func codeWalk(files []string,
	fileTypes []string,
	delayTimeInMsChannel chan time.Duration,
	colorChannel chan bool,
	snapShotChannel chan bool) {
	var currentDelay time.Duration = 2000000

	fileMap := make(map[string][]string)

	for _, file := range files {
		var extension = filepath.Ext(file)
		if contains(fileTypes, extension) {
			content, err := readFile(file)
			if err == nil {
				fileMap[file] = content
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

	for fileName, content := range fileMap {
		Info.Println("Current File:", fileName)
		for _, l := range content {
			for _, x := range l {

				select {
				case delayTimeInMs := <-delayTimeInMsChannel:
					currentDelay = delayTimeInMs
				default:
				}

				select {
				case snapShotSignal := <-snapShotChannel:
					if snapShotSignal {
						Info.Println("SnapShot current file: ", fileName)
					}
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

				time.Sleep(currentDelay * time.Nanosecond)
				fmt.Print(string(x))
			}
			fmt.Println("")
		}
	}
}

func keyHandler(delayChannel chan time.Duration, colorChannel chan bool, snapShotChannel chan bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	var delay time.Duration = 200
	var delayStep time.Duration = DELAY_STEP

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}
			if ev.Ch == '+' {
				delay, delayStep = decreaseDelay(delay, delayStep, delayChannel)
			}
			if ev.Ch == '-' {
				delay, delayStep = increaseDelay(delay, delayStep, delayChannel)
			}
			if ev.Ch == 'c' {
				colorChannel <- true
			}
			if ev.Ch == 's' {
				snapShotChannel <- true
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			break mainloop
		}
	}
}

func decreaseDelay(delay time.Duration,
	delayStep time.Duration,
	delayChannel chan time.Duration) (time.Duration, time.Duration) {
	if delay > 0+delayStep {
		delay = delay - delayStep
	} else {
		delayStep = delayStep / 2
		if delay > 0+delayStep {
			delay = delay - delayStep
		}
	}
	Info.Println("Decreased Current Delay:", delay)
	Info.Println("Adapt DelayStep to: ", delayStep)
	delayChannel <- delay
	return delay, delayStep
}

func increaseDelay(delay time.Duration,
	delayStep time.Duration,
	delayChannel chan time.Duration) (time.Duration, time.Duration) {
	delayStep = delayStep * 2
	delay = delay + delayStep
	Info.Println("Increased Current Delay:", delay)
	Info.Println("Adapt DelayStep to: ", delayStep)
	delayChannel <- delay
	return delay, delayStep
}

func main() {
	initLogging(DEFAULT_LOG)
	delayChannel := make(chan time.Duration)
	colorChannel := make(chan bool)
	snapShotChannel := make (chan bool)
	var fileTypes = []string{".tf", ".sh", ".java", ".go"}
	dir := flag.String("dir", "/", "directory to walk")
	flag.Parse()
	Info.Println("Walking Directory: ", *dir)
	var files []string
	err := filepath.Walk(*dir, visit(&files))
	if err != nil {
		panic(err)
	}
	Info.Println("Loaded ", len(files), "files")
	go codeWalk(files, fileTypes, delayChannel, colorChannel, snapShotChannel)
	keyHandler(delayChannel, colorChannel, snapShotChannel)
}
