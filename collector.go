package main

import (
	"os/exec"
	"strings"
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

func getCommitDates(path string) (string, string ,error) {
	fileName, dir := cutFileNameFromPath(path)
	cmdString := "cd " + dir + "&& git log " + fileName + " | grep Date"
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	output, err := cmd.Output()
	datesRaw := string(output)
	dates := strings.Split(datesRaw, "\n")
	firstCommitDate := "unknown"
	lastCommitDate := "unknown"
	if len(dates) > 1 {
		lastCommitDate = dates[0]
		firstCommitDate = dates[len(dates)-2]
		firstCommitDate = removeWordFromString(firstCommitDate, "Date:")
		lastCommitDate = removeWordFromString(lastCommitDate, "Date:")
	}
	return firstCommitDate, lastCommitDate, err
}
