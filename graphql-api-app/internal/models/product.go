package models

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Availability bool   `json:"availability"`
	Comment     string  `json:"comment"`
}