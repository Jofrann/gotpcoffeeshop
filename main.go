package main

import (
	"encoding/json"
	"fmt"
	"myapi/HTTP-Serveur/Go_API-Coffee-Shop/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var drinks []models.Drink
var orders []models.Order
var orderCounter int = 1

func generateOrderID() string {
	id := fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	return id
}

// GET /menu
func getMenuHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drinks)
}

// GET /menu/{id}
func getDrinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	for _, drink := range drinks {
		if drink.ID == id {
			json.NewEncoder(w).Encode(drink)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Boisson non trouvée"})
}

// POST /orders
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	// Boisson ?
	var selectedDrink *models.Drink
	for _, d := range drinks {
		if d.ID == order.DrinkID {
			selectedDrink = &d
			break
		}
	}

	if selectedDrink == nil {
		http.Error(w, "Boisson introuvable", http.StatusBadRequest)
		return
	}

	order.ID = generateOrderID()
	order.DrinkName = selectedDrink.Name
	order.OrderedAt = time.Now()
	order.Status = models.StatusPending
	order.TotalPrice = calculatePrice(selectedDrink.BasePrice, order.Size, order.Extras)

	orders = append(orders, order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GET /orders
func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GET /orders/{id}
func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	for _, order := range orders {
		if order.ID == id {
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
}

// PATCH /orders/{id}/status
func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = models.OrderStatus(payload.Status)
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
}

// DELETE /orders/{id}
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, order := range orders {
		if order.ID == id {

			if order.Status == models.StatusPickedUp {
				http.Error(w, "Impossible : commande déjà récupérée", http.StatusBadRequest)
				return
			}

			orders = append(orders[:i], orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Commande introuvable", http.StatusNotFound)
}

func calculatePrice(basePrice float64, size string, extras []string) float64 {
	price := basePrice

	switch size {
	case "small":
		price *= 0.8
	case "medium":
		price *= 1.0
	case "large":
		price *= 1.3
	}

	price += float64(len(extras)) * 0.50

	return price
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	drinks = []models.Drink{
		{ID: "DRK-001", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
		{ID: "DRK-002", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
		{ID: "DRK-003", Name: "Latte", Category: "coffee", BasePrice: 4.00},
		{ID: "DRK-004", Name: "Americano", Category: "coffee", BasePrice: 3.00},
		{ID: "DRK-005", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
		{ID: "DRK-006", Name: "Iced Coffee", Category: "coffee", BasePrice: 4.50},
	}

	r := mux.NewRouter()

	r.HandleFunc("/menu", getMenuHandler).Methods("GET")
	r.HandleFunc("/menu/{id}", getDrinkHandler).Methods("GET")

	r.HandleFunc("/orders", getOrdersHandler).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/status", updateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenue dans l'API Coffee Shop ☕"))
	})

	fmt.Println("🚀 Serveur lancé sur : http://localhost:8080")
	http.ListenAndServe(":8080", cors(r))
}
