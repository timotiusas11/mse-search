package interfaces

import (
	"encoding/json"
	"mse-search/usecases"
	"net/http"
)

type RecipeInteractor interface {
	SearchRecipe(keyword string, isHalal, isVegetarian bool, page, take int) (res []usecases.Recipe, err error)
}

type WebServiceHandler struct {
	RecipeInteractor RecipeInteractor
}

type SearchParam struct {
	Keyword      string `json:"keyword"`
	IsHalal      bool   `json:"isHalal"`
	IsVegetarian bool   `json:"isVegetarian"`
	Page         int    `json:"page"`
	Take         int    `json:"take"`
}

func (handler WebServiceHandler) SearchRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var req SearchParam

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipes, err := handler.RecipeInteractor.SearchRecipe(req.Keyword, req.IsHalal, req.IsVegetarian, req.Page, req.Take)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		result, err := json.Marshal(recipes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
