package parsers

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

// Checks the first letter of the word to determine of it is capitalized
// Using this to determine ingredient title
func isCapital(word string) bool {
	if len(word) == 0 {
		return false
	}
	r := rune(word[0])
	return unicode.IsUpper(r)
}

func ImportMarkdown(filepath string) (Recipe, error) {
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
				// Adding this to support seperating out ingredients i.e. salad and dressing for chicken salad
				if strings.HasPrefix(line, "-") {
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
				}
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
