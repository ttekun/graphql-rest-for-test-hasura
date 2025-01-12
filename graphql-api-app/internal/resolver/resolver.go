package resolver

import (
	"context"

	"graphql-api-app/internal/models"
)

type Resolver struct{}

func (r *Resolver) QueryProducts(ctx context.Context) ([]*models.Product, error) {
	// TODO: 製品を取得してProductスライスを返すロジックを実装する
	products := []*models.Product{
		{ID: "1", Name: "Product 1", Price: 9.99, Availability: true, Comment: "Great product"},
		{ID: "2", Name: "Product 2", Price: 24.99, Availability: false, Comment: "Out of stock"},
	}
	return products, nil
}
