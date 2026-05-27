package main

import (
	"bufio"
	"os"
	"regexp"
)

func main() {
	in, _ := os.Open("coverage.tmp")
	defer in.Close()

	out, _ := os.Create("coverage.out")
	defer out.Close()

	re := regexp.MustCompile(`mocks|main.go|api.go|config.go|test|_coverage.go|docs`)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if !re.MatchString(line) {
			out.WriteString(line + "\n")
		}
	}
}
