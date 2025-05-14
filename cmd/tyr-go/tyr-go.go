package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Ingredient struct {
	quantity    uint8    //keeping these small if you need over 255 of something thats a bad recipe
	measurement string   //Going to need to specify these
	name        string   //name of the ingredient
	preparation string   //how to prepare the ingredient
	meta        []string //metadata about the ingredient
}

type Recipe struct {
	title       string       //title of the recipe
	ingredients []Ingredient //slice of ingredient structures to store ingredients
	steps       []string     //list of steps for the recipe
}

var myData Recipe

func importMarkdown(filepath string) error {
	// Open the file, check for error
	file, err := os.Open(filepath)
	if err != nil {
		// TODO: improve error handling
		return err
	}
	// Close file after this function finishes
	defer file.Close()

	// Open a scanner and create the recipe
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Grab first line and make sure
		line := scanner.Text()

		// Check for leading # character, remember the first returned value is the index so _ throws it away
		for _, ch := range line {
			if unicode.IsSpace(ch) {
				continue // Skip any leading whitespace
			}
			if ch == '#' {
				//Recipe.title := line[1:]
			}
		}

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
	//for x, y := range myData {
	//	fmt.Println(x)
	//	fmt.Println(y)
	//}
}
