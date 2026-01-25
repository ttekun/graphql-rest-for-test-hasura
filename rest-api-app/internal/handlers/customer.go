package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ttekun/rest-api-app/internal/models"
)

// GetCustomerInfo handles the GET request for customer information
func GetCustomerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customers := models.GetMockCustomers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}