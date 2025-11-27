package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//
// ────────────────────────────────────────────────
//   STRUCTURES & DONNÉES
// ────────────────────────────────────────────────
//

// Drink représente une boisson du menu
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"` // coffee, tea, cold
	BasePrice float64 `json:"base_price"`
}

// OrderStatus représente l'état d'une commande
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"	
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked-up"	
	StatusCancelled OrderStatus = "cancelled"
)

// Order représente une commande
type Order struct {
	ID           string      `json:"id"`
	DrinkID      string      `json:"drink_id"`
	DrinkName    string      `json:"drink_name"`
	Size         string      `json:"size"`
	Extras       []string    `json:"extras"`
	CustomerName string      `json:"customer_name"`
	Status       OrderStatus `json:"status"`
	TotalPrice   float64     `json:"total_price"`
	OrderedAt    time.Time   `json:"ordered_at"`
}

// Base mémoire
var drinks []Drink
var orders []Order
var orderCounter int = 1

//
// ────────────────────────────────────────────────
//   HANDLERS MENU
// ────────────────────────────────────────────────
//

// GET /menu
func getMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drinks)
}

// GET /menu/{id}
func getDrink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	for _, d := range drinks {
		if d.ID == id {
			json.NewEncoder(w).Encode(d)
			return
		}
	}

	http.Error(w, "Boisson introuvable", http.StatusNotFound)
}

//
// ────────────────────────────────────────────────
//   HANDLERS COMMANDES
// ────────────────────────────────────────────────
//

// POST /orders
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Body JSON invalide", http.StatusBadRequest)
		return
	}

	// Vérifier boisson
	var drink *Drink
	for _, d := range drinks {
		if d.ID == order.DrinkID {
			drink = &d
			break
		}
	}
	if drink == nil {
		http.Error(w, "Boisson introuvable", http.StatusBadRequest)
		return
	}

	// Remplir
	order.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	order.DrinkName = drink.Name
	order.Status = StatusPending
	order.OrderedAt = time.Now()
	order.TotalPrice = calculatePrice(drink.BasePrice, order.Size, order.Extras)

	orders = append(orders, order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GET /orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GET /orders/{id}
func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	for _, o := range orders {
		if o.ID == id {
			json.NewEncoder(w).Encode(o)
			return
		}
	}

	http.Error(w, "Commande introuvable", http.StatusNotFound)
}

// PATCH /orders/{id}/status
func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	for i, o := range orders {
		if o.ID == id {
			orders[i].Status = OrderStatus(payload.Status)
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}

	http.Error(w, "Commande introuvable", http.StatusNotFound)
}

// DELETE /orders/{id}
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for i, o := range orders {
		if o.ID == id {

			if o.Status == StatusPickedUp {
				http.Error(w, "Commande déjà récupérée", http.StatusBadRequest)
				return
			}

			orders = append(orders[:i], orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Commande introuvable", http.StatusNotFound)
}

//
// ────────────────────────────────────────────────
//   PRIX
// ────────────────────────────────────────────────
//

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

//
// ────────────────────────────────────────────────
//   CORS
// ────────────────────────────────────────────────
//

func corsMiddleware(next http.Handler) http.Handler {
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

//
// ────────────────────────────────────────────────
//   MAIN
// ────────────────────────────────────────────────
//

func main() {

	drinks = []Drink{
		{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
		{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
		{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 4.00},
		{ID: "4", Name: "Americano", Category: "coffee", BasePrice: 3.00},
		{ID: "5", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
	}

	r := mux.NewRouter()

	r.HandleFunc("/menu", getMenu).Methods("GET")
	r.HandleFunc("/menu/{id}", getDrink).Methods("GET")

	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/status", updateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Coffee Shop ☕ — Ready"))
	})

	fmt.Println("🚀 Serveur sur http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware(r))
}
