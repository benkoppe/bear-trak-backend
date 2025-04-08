package external

type Location struct {
	Name   string `json:"location_name"`
	Number string `json:"location_number"`
}

type Event struct {
	LocationName   string `json:"location_name"`
	LocationNumber string `json:"location_number"`

	MealName   string `json:"meal_name"`
	MealNumber string `json:"meal_number"`

	MenuCategoryName   string `json:"menu_category_name"`
	MenuCategoryNumber string `json:"menu_category_number"`

	ServeDate string `json:"serve_date"`
}

type Recipe struct {
	ID int `json:"ID"`

	LocationName   string `json:"Location_Name"`
	LocationNumber string `json:"Location_Number"`

	MealName   string `json:"Meal_Name"`
	MealNumber int    `json:"Meal_Number"`

	MenuCategoryName   string `json:"Menu_Category_Name"`
	MenuCategoryNumber string `json:"Menu_Category_Number"`

	RecipeName   string `json:"Recipe_Print_As_Name"`
	RecipeNumber string `json:"Recipe_Number"`

	ServeDate string `json:"Serve_Date"`
}
