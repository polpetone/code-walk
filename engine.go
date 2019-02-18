package main

import (
	"github.com/fatih/color"
	"math/rand"
	"path/filepath"
	"time"
)

const DELAY_STEP = 100 * time.Millisecond

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

type CodeWalkFileInfo struct {
	ProjectName string
	Authors []string
	FileName string
	FirstCommitDate time.Time
	LastCommitDate time.Time
}

var delayChannel = make(chan time.Duration, 10)
var snapShotChannel = make(chan bool, 10)
var haltChannel = make(chan bool)
var jumpFileChannel = make(chan bool)
var continueChannel = make(chan bool)
var codeWalkFileInfoChannel = make(chan CodeWalkFileInfo)

func engine(files []string, fileTypes []string) {
	Info.Println("Loaded ", len(files), "files")
	go codeWalk(files, fileTypes)
}

func sendCodeWalkFileInfo(fileName string){
	firstCommitDate, lastCommitDate , err := getCommitDates(fileName)
	if err != nil {
		Error.Println(err)
	}
	authors, _ := getGitAuthors(fileName)
	fileNameWithoutWalkDir := removeWordFromString(fileName, directoryToWalk)
	projectName, _ := gitProjectNameAndAbsolutePathFromFilePath(fileName)
	codeWalkFileInfoChannel <-
		CodeWalkFileInfo{
			ProjectName : projectName,
			FileName:fileNameWithoutWalkDir,
			Authors: authors,
			FirstCommitDate: firstCommitDate,
			LastCommitDate: lastCommitDate,
		}
}

//TODO: optimize: clear screen when jump to next file. jump immediately,
//currently it seems to wait until complete line is print
func codePrinter(tactChannel chan bool, snapShotChannel chan bool, jumpFileChannel chan bool, contentMap map[string][]string) {
	keys := keysFromMap(contentMap)
	jumpFile := false
	for iterator := 0; iterator < len(contentMap); iterator++ {
			fileName := keys[iterator]
			content := contentMap[fileName]
			Trace.Println(fileName)
			Trace.Println("Show File Number",iterator+1, "from total loaded files", len(contentMap))
			if jumpFile {
				jumpFile = false
				Trace.Println("jump file")
				continue
			}
			go sendCodeWalkFileInfo(fileName)
			for _, l := range content {

				if jumpFile {
					Trace.Println("jump line")
					break
				}
				for _, x := range l {
					select {
					case <-tactChannel:
						codeChannel <- string(x)
					case <-jumpFileChannel:
						Trace.Println("Jump File received")
						jumpFile = true
						break
					case <-snapShotChannel:
						var snapShotFile= "snapshot-" + time.Now().Format("2006-01-02_15:04:05")
						Info.Println("SnapShot current file: ", fileName)
						err := appendLineToFile(codeWalkDir+"/"+snapShotFile, fileName)
						if err != nil {
							Error.Println("Failed to write snapshot: ", err)
						}
					}
				}
				codeChannel <- "\n"
			}
			if iterator == len(contentMap) - 1{
				Trace.Println("Reached End of loaded Files. Start from beginning again.")
				iterator = -1
			}
	}
}

func codeWalk(files []string,
	fileTypes []string) {

	tactChannel := make(chan bool)

	var halt = false
	var currentDelay = 50 * time.Millisecond
	fileContentMap := loadFileContentMap(files, fileTypes)
	color.Set(color.FgHiGreen)
	rand.Seed(time.Now().Unix())
	defer color.Unset()

	go codePrinter(tactChannel, snapShotChannel, jumpFileChannel, fileContentMap)

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
		case delayTimeInMs := <-delayChannel:
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
