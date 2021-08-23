package domain

type Recipe struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Ingredients  []string `json:"ingredients"`
	IsHalal      bool     `json:"is_halal"`
	IsVegetarian bool     `json:"is_vegetarian"`
	Description  string   `json:"description"`
	Rating       float64  `json:"rating"`
}

type RecipeRepository interface {
	Search(keyword string, isHalal, isVegetarian bool, page, take int) ([]Recipe, error)
}
