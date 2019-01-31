package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func getGitAuthors(path string) ([]string, error) {
	fileName, dir := cutFileNameFromPath(path)
	cmdString := "cd " + dir + "&& git blame " + fileName + " --show-stats -p | grep '^author ' | uniq"
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	output, err := cmd.Output()
	authorsRaw := string(output)
	authorsRaw2 := strings.Split(authorsRaw, "\n")
	var authors []string
	for _,a := range authorsRaw2 {
		b := removeWordFromString(a, "author")
		authors = append(authors, b)
	}
	return uniqueNonEmptyElementsOf(authors), err
}

func getCommitDates(path string) (time.Time, time.Time , error) {
	fileName, dir := cutFileNameFromPath(path)
	cmdString := "cd " + dir + "&& git log  --date=iso " + fileName + " | grep Date"
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	output, err := cmd.Output()
	datesRaw := string(output)
	dates := strings.Split(datesRaw, "\n")

	var firstTime time.Time
	var lastTime time.Time
	if len(dates) > 1 {
		lastCommitDate := dates[0]
		firstCommitDate := dates[len(dates)-2]

		firstCommitDate = removeWordFromString(firstCommitDate, "Date:")
		lastCommitDate = removeWordFromString(lastCommitDate, "Date:")

		firstTime, err = parseTime(firstCommitDate)
		lastTime, err = parseTime(lastCommitDate)
	}
	return firstTime, lastTime, err
}

//FIXME still a bug in some constellations
func parseTime(date string) (time.Time, error) {
	layout_negativ := "2006-01-02 15:04:05 -0100"
	layout_positiv := "2006-01-02 15:04:05 +0100"
	parsed, err := time.Parse(layout_positiv, date)
	if err == nil {
		return parsed, nil
	} else {
		fmt.Println("time parsing failed layout positiv, trying layout negativ")
		return time.Parse(layout_negativ, date)
	}
}
