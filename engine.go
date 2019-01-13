package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
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

func engine(files []string, fileTypes []string) {
	delayChannel := make(chan time.Duration, 10)
	snapShotChannel := make(chan bool, 10)
	haltChannel := make(chan bool)
	continueChannel := make(chan bool)

	Info.Println("Loaded ", len(files), "files")

	go codeWalk(files, fileTypes, delayChannel, snapShotChannel, haltChannel, continueChannel)
	keyHandler(delayChannel, snapShotChannel, haltChannel, continueChannel)
}

func codePrinter(tactChannel chan bool, snapShotChannel chan bool, contentMap map[string][]string) {
	for fileName, content := range contentMap {
		go client(fileName)
		for _, l := range content {
			for _, x := range l {
				select {
				case <-tactChannel:
					fmt.Print(string(x))
				case <-snapShotChannel:
					var snapShotFile = "snapshot-" + time.Now().Format("2006-01-02_15:04:05")
					Info.Println("SnapShot current file: ", fileName)
					err := appendLineToFile(codeWalkDir+"/"+snapShotFile, fileName)
					if err != nil {
						Error.Println("Failed to write snapshot: ", err)
					}
				}
			}
			fmt.Println("")
		}
	}
}

func codeWalk(files []string,
	fileTypes []string,
	delayTimeInMsChannel chan time.Duration,
	snapShotChannel chan bool,
	haltChannel chan bool,
	continueChannel chan bool) {

	tactChannel := make(chan bool)

	var halt = false
	var currentDelay = 100 * time.Millisecond
	fileContentMap := loadFileContentMap(files, fileTypes)
	color.Set(color.FgHiGreen)
	rand.Seed(time.Now().Unix())
	defer color.Unset()

	go codePrinter(tactChannel, snapShotChannel, fileContentMap)

	for {
		time.Sleep(currentDelay)

		if !halt {
			tactChannel <- true
		}

		select {
		case colorSwitch := <-colorChannel:
			randomColorIndex := rand.Intn(6)
			if colorSwitch {
				color.Set(colors[randomColorIndex])
			}
		case delayTimeInMs := <-delayTimeInMsChannel:
			currentDelay = delayTimeInMs
		case <-haltChannel:
			halt = true
		case <- continueChannel:
			halt = false
		default:
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
