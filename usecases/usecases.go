package usecases

import (
	"mse-search/domain"
)

type Recipe struct {
	Id           int
	Name         string
	Ingredients  []string
	IsHalal      bool
	IsVegetarian bool
	Description  string
	Rating       float64
}

type RecipeInteractor struct {
	RecipeRepository domain.RecipeRepository
}

func (interactor *RecipeInteractor) SearchRecipe(keyword string, isHalal, isVegetarian bool, page, take int) (res []Recipe, err error) {
	recipes, err := interactor.RecipeRepository.Search(keyword, isHalal, isVegetarian, page, take)
	if err != nil {
		return nil, err
	}

	for _, rec := range recipes {
		res = append(res, Recipe{
			Id:           rec.Id,
			Name:         rec.Name,
			Ingredients:  rec.Ingredients,
			IsHalal:      rec.IsHalal,
			IsVegetarian: rec.IsVegetarian,
			Description:  rec.Description,
			Rating:       rec.Rating,
		})
	}
	return res, nil
}
