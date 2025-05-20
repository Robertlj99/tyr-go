package parsers

// Structure to hold ingredients fields defined as follows
//
//	Quantity: How much of the ingredient is to be used, may convert these to int during
//			  computations with it but keeping string to represent things like 1/4
//
//	Measurement: Bit misleading, this keeps the measurement label i.e.; cups, tsp, etc.
//
//	Name: Not misleading, stores the name
//
//	Preparation: Stores how to prep the ingredient, i.e.; chopped, shredded, etc
type Ingredient struct {
	Quantity    string
	Measurement string
	Name        string
	Preparation string
}

// Structure to hold Recipes fields defined as follows
//
//	Title: Title of the recipe
//
//	Ingredients: A slice of ingredients to be used in the recipe
//
//	Steps: A slice of strings explaining how to prepare the ingredients (aka Instructions)
type Recipe struct {
	Title       string
	Ingredients []Ingredient
	Steps       []string
	Category    string
}
