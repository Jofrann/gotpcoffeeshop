package handlers

import (
    "encoding/json"
    "net/http"

    "gotpcoffeeshop/models"
)

var Drinks []models.Drink

// GET /menu
func GetMenu(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Drinks)
}

// GET /menu/{id}
func GetDrink(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.PathValue("id")

    for _, drink := range Drinks {
        if drink.ID == id {
            json.NewEncoder(w).Encode(drink)
            return
        }
    }

    http.Error(w, "Boisson non trouvée", http.StatusNotFound)
}
