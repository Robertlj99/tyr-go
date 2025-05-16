package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Robertlj99/tyr-go/internal/parsers"
)

func main() {
	// Get current working directory
	// This is used to get the path to the recipe file
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory")
		panic(err)
	}
	fmt.Println("\nCurrent working directory: ", cwd)

	// Go through directory tree until you arrive at top tyr-go directory
	for {
		base := filepath.Base(cwd)          // Name of current directory
		parent := filepath.Dir(cwd)         // Full path of parent directory
		parentBase := filepath.Base(parent) // Name of parent directory

		// Check to make sure we are in tyr-go directory, double check to make sure it is not
		// the tyr-go/cmd/tyr-go/ directory, we want to get out of that and into just tyr-go/
		if base == "tyr-go" && parentBase != "cmd" {
			fmt.Println("Found target directory", cwd)
			break
		}

		// If we can find the root tyr-go/ directory then the rest of the files will not be where
		// the program expects them so it is not time to panic
		if parent == cwd {
			fmt.Println("Reached root without finding valid directory, please reinstall or repair")
			return
		}

		// If we get here we did not panic but we are not in the right directory so move up the chain
		cwd = parent
	}

	//EXISTING WORK
	// Create the path to the recipe directory
	filepath := filepath.Join(cwd, "assets", "md-recipes")
	fmt.Printf("File path: %s \n \n", filepath)

	// Open the file and parse it
	recipe, err := parsers.ImportMarkdown(filepath)
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
	fmt.Println()
}
