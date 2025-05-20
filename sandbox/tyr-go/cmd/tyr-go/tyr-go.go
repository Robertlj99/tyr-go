package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

	// Finally, make the filepath
	dirPath := filepath.Join(cwd, "assets", "md-recipes")
	fmt.Printf("File path: %s \n \n", dirPath)

	// Slice to store the recipes
	var recipes []parsers.Recipe

	// Tracking categories, may remove this but may be useful
	categories := make(map[string]int)

	// Recursive walk through directories
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Keep trying to access files despite error
		}

		// Don't process directories
		if info.IsDir() {
			return nil
		}

		// Don't process any non-markdown files
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Import the recipe
		recipe, err := parsers.ImportMarkdown(path)
		if err != nil {
			fmt.Printf("Error importing recipe %s: %v\n", path, err)
			return nil
		}

		//Get category using the path
		relPath, _ := filepath.Rel(dirPath, path)
		category := filepath.Dir(relPath)
		if category == "." {
			category = "uncategorized"
		}

		categories[category]++
		recipe.Category = category

		recipes = append(recipes, recipe)

		return nil
	})

	if err != nil {
		fmt.Printf("Error while walking filepath: %v\n", err)
		return
	}

	fmt.Printf("\nSuccesfull import of %d recipes\n", len(recipes))
	fmt.Println("\nRecipes by category:")
	for category, count := range categories {
		fmt.Printf("- %s: %d recipes\n", category, count)
	}

	fmt.Println("\nIndividiual Recipes")
	for i, recipe := range recipes {
		fmt.Printf("\n%d. %s\n", i+1, recipe.Title)
		for _, ingredients := range recipe.Ingredients {
			fmt.Printf("Quantity: %-10s Measurement: %-10s Name: %-30s Prep: %-10s \n",
				ingredients.Quantity, ingredients.Measurement, ingredients.Name, ingredients.Preparation)
		}
		fmt.Println("\nSteps:")
		for _, steps := range recipe.Steps {
			fmt.Println(steps)
		}
	}
}
