package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"graphql-api-app/internal/resolver"

	"github.com/graphql-go/graphql"
)

// Product represents a product in the system
type Product struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Availability bool    `json:"availability"`
	Comment      string  `json:"comment"`
}

func main() {
	resolver := &resolver.Resolver{}

	fields := graphql.Fields{
		"queryProducts": &graphql.Field{
			// Return the list of products
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
				// Call the QueryProducts function from resolver.go
				products, err := resolver.QueryProducts(context.Background())
				if err != nil {
					return nil, err
				}
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

	// Define the GraphQL handler with enhanced configuration
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			// Fallback to URL query for GET requests
			params.Query = r.URL.Query().Get("query")
		}

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: params.Query,
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
