package main

import (
	"fmt"
	"log"
	"mse-search/infrastructure"
	"mse-search/interfaces"
	"mse-search/usecases"
	"net/http"
)

func main() {
	client, err := infrastructure.GetESClient()
	if err != nil {
		return
	}
	log.Println("ES initialized.")

	repo := interfaces.NewRecipeRepository(client)
	recipeInteractor := new(usecases.RecipeInteractor)
	recipeInteractor.RecipeRepository = repo

	webServiceHandler := interfaces.WebServiceHandler{}
	webServiceHandler.RecipeInteractor = recipeInteractor

	http.HandleFunc("/search", webServiceHandler.SearchRecipeHandler)
	http.HandleFunc("/autocomplete", webServiceHandler.SearchRecipeHandler)

	fmt.Println("Starting Web Server at http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
