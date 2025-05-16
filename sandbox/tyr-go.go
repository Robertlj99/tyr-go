package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	var current uint8 // Int to determine which section it is: 1 = ingredients, 2 = steps
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
				// Find title
				line = strings.TrimPrefix(line, "-")
				parts := strings.Fields(line)
				nameStart := -1

				// Iterate through each string (parts) word by word (word)
				for i, word := range parts {
					if isCapital(word) {
						nameStart = i
						break
					}
				}

				// If no capitalized word was found add empty ingredient
				if nameStart == -1 {
					ingredient := Ingredient{Name: "(Error)"}
					recipe.Ingredients = append(recipe.Ingredients, ingredient)
					break
				}

				// Find end of the name (last consecutive capitalized word)
				nameEnd := nameStart

				// Iterate from the first capitalized word until the line ends
				for i := nameStart; i < len(parts); i++ {
					if isCapital(parts[i]) {
						nameEnd = i
					}
				}

				// Build the name
				nameSlice := parts[nameStart : nameEnd+1]
				name := strings.Join(nameSlice, " ")

				// Split rest of string into before and after
				before := parts[:nameStart]
				after := parts[nameEnd+1:]

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
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
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
	// Get current working directory
	// This is used to get the path to the recipe file
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory")
		panic(err)
	}
	fmt.Println("\nCurrent working directory: ", cwd)

	// Get the parent directory of the current working directory
	parentDir := filepath.Dir(cwd)
	fmt.Println("Parent directory: ", parentDir)

	// Create the path to the recipe file
	filepath := filepath.Join(parentDir, "assets", "md-recipes", "Creamy Chicken Enchilada Soup.md")
	fmt.Printf("File path: %s \n \n", filepath)

	// Open the file and parse it
	recipe, err := importMarkdown(filepath)
	if err != nil {
		fmt.Printf("Scanner errored: %e", err)
	}

	// Print to test here
	fmt.Printf("Recipe Title: %s \n", recipe.Title)
	ingredients := recipe.Ingredients
	steps := recipe.Steps
	fmt.Printf("\nNumber of ingredients: %v \n", len(ingredients))
	fmt.Println("Listing out ingredients")
	for i := range len(ingredients) {
		fmt.Printf("Quantity: %-10s Measurement: %-10s Name: %-30s Prep: %-10s \n",
			ingredients[i].Quantity, ingredients[i].Measurement, ingredients[i].Name, ingredients[i].Preparation)
	}
	fmt.Println("\nNumber of steps: ", len(steps))
	fmt.Println("Printing Steps")
	for i := range steps {
		fmt.Println(steps[i])
	}
}
