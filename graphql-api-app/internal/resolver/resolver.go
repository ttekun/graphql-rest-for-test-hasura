package main

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) QueryProducts(ctx context.Context) ([]*Product, error) {
	// TODO: Implement the logic to query and return products
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Availability bool   `json:"availability"`
	Comment     string  `json:"comment"`
}