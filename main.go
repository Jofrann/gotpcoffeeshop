package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Go-TP-API-Coffee-Shop-avec-Mux-/models"
	"github.com/gorilla/mux"
)


// === Structures ===

// === Données globales ===
var drinks []models.Drink
var orders []models.Order
var orderCounter int = 1

func generateOrderID() string {
	id := fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	return id
}

// === Handlers ===
func getMenuHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drinks)
}

func getDrinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Récupérer la variable dans l’URL
	vars := mux.Vars(r)
	id := vars["id"] // "DRK-001", "DRK-002", etc.

	// 2. Parcourir la liste des boissons
	for _, drink := range drinks {
		if drink.ID == id {
			// 3. Si on trouve la bonne boisson
			json.NewEncoder(w).Encode(drink)
			return
		}
	}

	// 4. Si on n’a rien trouvé → 404
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Boisson non trouvée"})
}

// POST /orders - Créer une nouvelle commande
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Corps JSON invalide", http.StatusBadRequest)
		return
	}

	// Vérifier que la boisson existe
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

	// Générer un ID unique
	order.ID = generateOrderID()

	// Compléter les champs
	order.DrinkName = selectedDrink.Name
	/*order.Status = StatusPending*/
	order.OrderedAt = time.Now()
	order.TotalPrice = calculatePrice(*selectedDrink, order.Size, order.Extras)

	// Enregistrer la commande
	orders = append(orders, order)

	// Répondre
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer_name := vars["id"]
	postId := vars["postId"]
	fmt.Fprintf(w, "Drinks %s, Post %s", customer_name, postId)
}

// === Middleware pour l'interface Web de test ===
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func calculatePrice(drink models.Drink, size string, extras []string) float64 {
	price := drink.BasePrice

	switch size {
	case "medium":
		price += 0.5
	case "large":
		price += 1.0
	}

	for range extras {
		price += 0.3
	}

	return price
}

// === Point d’entrée ===
func main() {
	// Initialisation des données
	drinks = []models.Drink{
		{ID: "DRK-001", Name: "Cappuccino", Category: "coffee", BasePrice: 2.0},
		{ID: "DRK-002", Name: "Hot Chocolate Deluxe", Category: "chocolate", BasePrice: 3.5},
		{ID: "DRK-003", Name: "Latte", Category: "coffee", BasePrice: 2.8},
	}

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/menu", getMenuHandler).Methods("GET")
	r.HandleFunc("/menu/{id}", getDrinkHandler).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")

	fmt.Println("🚀 Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware(r))

}
