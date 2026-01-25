package models

// Customer represents a customer in the system
type Customer struct {
	LoginID string `json:"loginId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// GetMockCustomers returns a list of mock customer data
func GetMockCustomers() []Customer {
	return []Customer{
		{
			LoginID: "john123",
			Name:    "John Doe",
			Email:   "john.doe@example.com",
			Address: "123 Main St",
		},
		{
			LoginID: "jane456",
			Name:    "Jane Smith",
			Email:   "jane.smith@example.com",
			Address: "456 Elm St",
		},
	}
}