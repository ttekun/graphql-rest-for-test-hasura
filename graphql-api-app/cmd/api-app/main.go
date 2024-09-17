package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Product struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Availability bool    `json:"availability"`
	Comment      string  `json:"comment"`
}

var products []Product

func init() {
	// Initialize dummy products
	products = []Product{
		{
			ID:           "1",
			Name:         "Product 1",
			Price:        9.99,
			Availability: true,
			Comment:      "This is product 1",
		},
		{
			ID:           "2",
			Name:         "Product 2",
			Price:        19.99,
			Availability: false,
			Comment:      "This is product 2",
		},
	}
}

func main() {
	// Define the GraphQL schema
	fields := graphql.Fields{
		"queryProducts": &graphql.Field{
			// productsを返す
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Product",
				Fields: graphql.Fields{
					"id":           &graphql.Field{Type: graphql.String},
					"name":         &graphql.Field{Type: graphql.String},
					"price":        &graphql.Field{Type: graphql.Float},
					"availability": &graphql.Field{Type: graphql.Boolean},
					"comment":      &graphql.Field{Type: graphql.String},
				},
			})),
			Description: "Get the list of products",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return products, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	// Define the GraphQL handler
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		if len(result.Errors) > 0 {
			log.Printf("GraphQL query error: %v", result.Errors)
		}
		json.NewEncoder(w).Encode(result)
	})

	// Start the server
	fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
