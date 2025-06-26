package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	var inputFiles []string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter log file paths (one per line). Press Enter on an empty line to finish:")

	for {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		inputFiles = append(inputFiles, line)
	}

	if len(inputFiles) == 0 {
		fmt.Println("No input files given.")
		return
	}

	outputFile := "errors.log"
	err := ProcessLogs(inputFiles, outputFile)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Done. Errors written to", outputFile)
	}
}
func ProcessLogs(inputFiles []string, outputFile string) error {
	var wg sync.WaitGroup
	errorChannel := make(chan string)
	go func() {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Cannot create output file:", err)
			return
		}
		defer f.Close()
		writer := bufio.NewWriter(f)

		for line := range errorChannel {
			writer.WriteString(line + "\n")
		}
		writer.Flush()
	}()
	for _, file := range inputFiles {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			f, err := os.Open(filename)
			if err != nil {
				fmt.Println("Cannot open file:", filename, err)
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "ERROR") {
					errorChannel <- line
				}
			}
		}(file)
	}
	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	return nil
}
