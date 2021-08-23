package interfaces

import (
	"context"
	"encoding/json"
	"mse-search/domain"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type QueryRoot struct {
	Query interface{} `json:"query,omitempty"`
}

type MatchQuery struct {
	Query string `json:"query,omitempty"`
}

type BoolQuery struct {
	Bool BoolQueryParams `json:"bool"`
}

type BoolQueryParams struct {
	Must               interface{} `json:"must,omitempty"`
	Should             interface{} `json:"should,omitempty"`
	Filter             interface{} `json:"filter,omitempty"`
	MinimumShouldMatch int         `json:"minimum_should_match,omitempty"`
}

type SearchResult struct {
	Hits HitsSearchResult `json:"hits"`
}

type HitsSearchResult struct {
	ArrayHits []ArrayHits `json:"hits"`
}

type ArrayHits struct {
	Source map[string]interface{} `json:"_source"`
}

type recipeRepo struct {
	client *elasticsearch.Client
}

func NewRecipeRepository(client *elasticsearch.Client) domain.RecipeRepository {
	return &recipeRepo{client: client}
}

func (repo *recipeRepo) Search(keyword string, isHalal, isVegetarian bool, page, take int) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	var data SearchResult

	dataJson, err := json.Marshal(
		map[string]interface{}{
			"from": page,
			"size": take,
			"query": map[string]interface{}{
				"function_score": map[string]interface{}{
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"filter": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"is_halal": isHalal,
									},
								},
								{
									"term": map[string]interface{}{
										"is_vegetarian": isVegetarian,
									},
								},
							},
							"must": []map[string]interface{}{
								{
									"match": map[string]interface{}{
										"name": map[string]interface{}{
											"query":     keyword,
											"fuzziness": "AUTO",
											"operator":  "and",
										},
									},
								},
							},
						},
					},
					"functions": []map[string]interface{}{
						{
							"field_value_factor": map[string]interface{}{
								"field":  "rating",
								"factor": 1.2,
							},
						},
					},
					"boost_mode": "multiply",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	js := string(dataJson)

	request := esapi.SearchRequest{
		Index: []string{"recipes"},
		Body:  strings.NewReader(js),
	}

	res, err := request.Do(context.Background(), repo.client)

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	for _, d := range data.Hits.ArrayHits {
		var ingredients []string

		for _, ing := range d.Source["ingredients"].([]interface{}) {
			ingredients = append(ingredients, ing.(string))
		}

		recipes = append(recipes, domain.Recipe{
			Id:           int(d.Source["id"].(float64)),
			Name:         d.Source["name"].(string),
			Ingredients:  ingredients,
			IsHalal:      d.Source["is_halal"].(bool),
			IsVegetarian: d.Source["is_vegetarian"].(bool),
			Description:  d.Source["description"].(string),
			Rating:       d.Source["rating"].(float64),
		})
	}

	return recipes, nil
}
