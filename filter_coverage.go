package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func main() {
	in, _ := os.Open("coverage.tmp")
	defer in.Close()

	out, _ := os.Create("coverage.out")
	defer out.Close()

	var ignoreList = []string{
		"mocks",
		"pkg",
		"utils",
		"main\\.go",
		"api\\.go",
		"config\\.go",
		"test",
		"_coverage\\.go",
		"docs",
	}
	re := regexp.MustCompile(strings.Join(ignoreList, "|"))

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if !re.MatchString(line) {
			out.WriteString(line + "\n")
		}
	}
}
