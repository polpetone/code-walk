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

	return authors, err
}
