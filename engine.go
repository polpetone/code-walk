package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"path/filepath"
	"time"
)

const DELAY_STEP = 20000

var colorChannel = make(chan bool)

var colors = []color.Attribute{
	color.FgHiGreen,
	color.FgBlue,
	color.FgCyan,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiWhite,
}

type Command int

func (c *Command) Receive(key string, reply *string) error {
	Info.Println("received key", key)
	colorChannel <- true
	return nil
}

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

func engine(files []string, fileTypes []string) {
	delayChannel := make(chan time.Duration)
	snapShotChannel := make(chan bool)
	haltChannel := make(chan bool)

	Info.Println("Loaded ", len(files), "files")

	go server()
	go codeWalk(files, fileTypes, delayChannel, snapShotChannel, haltChannel)
	keyHandler(delayChannel, snapShotChannel, haltChannel)
}

func keyHandler(
	delayChannel chan time.Duration,
	snapShotChannel chan bool,
	haltChannel chan bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	var delay time.Duration = 200
	var delayStep time.Duration = DELAY_STEP

	defer Info.Println("KeyHandler stopped")

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}
			Info.Println("Pressed: ", ev.Ch)
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
			if ev.Ch == 'h' {
				haltChannel <- true
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			break mainloop
		}
	}
}

func loadFileContentMap(files []string, fileTypes []string) map[string][]string {
	fileContentMap := make(map[string][]string)
	for _, file := range files {
		var extension = filepath.Ext(file)
		if contains(fileTypes, extension) {
			content, err := readFile(file)
			if err == nil {
				fileContentMap[file] = content
			}
		}
	}
	return fileContentMap
}

func codeWalk(files []string,
	fileTypes []string,
	delayTimeInMsChannel chan time.Duration,
	snapShotChannel chan bool,
	haltChannel chan bool) {

	var currentDelay time.Duration = 2000000

	fileContentMap := loadFileContentMap(files, fileTypes)

	color.Set(color.FgHiGreen)

	rand.Seed(time.Now().Unix())

	defer color.Unset()

	var snapShotFile = "snapshot-" + time.Now().Format("2006-01-02_15:04:05")

	var haltSwitch = false
	var lastDelay time.Duration

	for fileName, content := range fileContentMap {
		for _, l := range content {
			for _, x := range l {

				select {
				case delayTimeInMs := <-delayTimeInMsChannel:
					currentDelay = delayTimeInMs
				default:
				}

				select {
				case <-haltChannel:
					Info.Println("Halt pressed")
					if haltSwitch ==  false{
						haltSwitch = true
						lastDelay = currentDelay
						currentDelay = 10 * time.Second
						Info.Println("Enabled Halt -> Current Delay: ", currentDelay)
					} else {
						haltSwitch = false
						currentDelay = lastDelay
						Info.Println("Disabled Halt -> Current Delay: ", currentDelay)
					}
				default:
				}

				select {
				case snapShotSignal := <-snapShotChannel:
					if snapShotSignal {
						Info.Println("SnapShot current file: ", fileName)
						err := appendLineToFile(codeWalkDir+"/"+snapShotFile, fileName)
						if err != nil {
							Error.Println("Failed to write snapshot: ", err)
						}
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
