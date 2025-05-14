package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Ingredient struct {
	Quantity    string //string to represent numbers like 1/4
	Measurement string //Going to need to specify these
	Name        string //name of the ingredient
	Preparation string //how to prepare the ingredient
	// meta        []string //metadata about the ingredient
}

type Recipe struct {
	Title       string
	Ingredients []Ingredient
	Steps       []string
}

// Checks the first letter of the word to determine of it is capitalized
// Using this to determine ingredient title
func isCapital(word string) bool {
	if len(word) == 0 {
		return false
	}
	r := rune(word[0])
	return unicode.IsUpper(r)
}

// Pass scanner from import markdown to this to get list of ingredients finish tommorow
// Current Strategy: Name is always uppercase so find that first, anything before that is quantity
// or measurement and anything afterwards is preparation instructions
func parseIngredients(scanner *bufio.Scanner) ([]Ingredient, error) {

	// Initialize slice of ingredients to store each one as its processed
	var ingredients []Ingredient

	// Scan through the scanner line by line processing each line as an ingredient
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Stop when you get to next section indicated by ##
		if strings.HasPrefix(line, "##") {
			return ingredients, nil
		}

		// All ingredients should start with '-', if not consider the line empty
		if strings.HasPrefix(line, "-") {

			// Find title
			parts := strings.Fields(line)
			index := -1
			var name string
			for i, word := range parts {
				if isCapital(word) {
					name = word
					index = i
					break
				}
			}

			// If no capitalized word was found add empty ingredient
			if index == -1 {
				ingredient := Ingredient{Name: "(Error)"}
				ingredients = append(ingredients, ingredient)
				break
			}

			// Split rest of string into before and after
			before := parts[:index]  // includes quantity and measurement if provided
			after := parts[index+1:] //includes prep instructions if provided

			// Grab the amount and measurement, if both provided length will be > 1, if only one
			// assume that it is amount, if none length will be 0
			var quantity, measurement string
			if len(before) > 0 {
				quantity = before[0]
			}
			if len(before) > 1 {
				measurement = before[1]
			}

			// Everything after title is prep instructions
			prep := strings.Join(after, " ")

			// Create ingredient and append it to ingredients list
			ingredient := Ingredient{
				Quantity:    quantity,
				Measurement: measurement,
				Name:        name,
				Preparation: prep,
			}
			ingredients = append(ingredients, ingredient)
		}
	}

	if err := scanner.Err(); err != nil {
		return ingredients, err
	}

	return ingredients, nil
}

func importMarkdown(filepath string) (Recipe, error) {
	var recipe Recipe

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return recipe, err
	}
	// Close file when this function closes
	defer file.Close()

	// Scan the text line by line
	scanner := bufio.NewScanner(file)
	// Int to determine which section it is: 1 = ingredients, 2 = steps
	var current uint8

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check section header against documented rules, if not sure check documentation
		// TODO: Include this in docs
		if strings.HasPrefix(line, "# ") {
			recipe.Title = strings.TrimPrefix(line, "# ")
		} else if strings.HasPrefix(line, "## Ingredients") {
			current = 1
		} else if strings.HasPrefix(line, "## Steps") || strings.HasPrefix(line, "## Instructions") {
			current = 2
		} else if line != "" {
			switch current {
			case 1:
				ingredients, err := parseIngredients(scanner)
				if err != nil {
					fmt.Printf("Error while parsing ingredients: %e", err)
				}
				recipe.Ingredients = append(recipe.Ingredients, ingredients...)
			case 2:
				recipe.Steps = append(recipe.Steps, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return recipe, err
	}

	return recipe, nil
}

func main() {
	fmt.Println("Hello")
	filepath := `C:\Users\robert.johnson\Desktop\VSCode\Git Repos\tyr-go\assets\md-recipes\Creamy Chicken Enchilada Soup.md`
	recipe, err := importMarkdown(filepath)
	if err != nil {
		fmt.Printf("Scanner errored: %e", err)
	}
	fmt.Printf("Recipe title is: %s\n", recipe.Title)
	fmt.Println("Ingredients as follows:")
	for x, y := range recipe.Ingredients {
		fmt.Printf("%v: %s", x, y)
		fmt.Println()
	}
	fmt.Println("Steps as follows:")
	for x, y := range recipe.Steps {
		fmt.Printf("%v: %s", x, y)
		fmt.Println()
	}

}
