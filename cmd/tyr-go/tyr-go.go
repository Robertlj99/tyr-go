package main

import (
	"bufio"
	"fmt"
	"os"
)

var myData []string

func importMarkdown(filepath string) error {
	// Open the file, check for error
	file, err := os.Open(filepath)
	if err != nil {
		// TODO: improve error handling
		return err
	}
	// Close file after this function finishes
	defer file.Close()

	// Scan text into myData line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		myData = append(myData, scanner.Text())
	}

	// Check if scanner errored
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Enter full path to document after calling executable")
		return
	}
	fmt.Println("Hello, Tyr!")
	err := importMarkdown(arguments[1])
	if err != nil {
		fmt.Println("Something went wrong")
	}
	for x := range myData {
		fmt.Println(x)
	}
}
