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

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}