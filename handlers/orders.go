package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "gotpcoffeeshop/models"
)

var Orders []models.Order
var OrderCounter = 1

func generateOrderID() string {
    id := fmt.Sprintf("ORD-%03d", OrderCounter)
    OrderCounter++
    return id
}

func CalculatePrice(base float64, size string, extras []string) float64 {
    price := base

    switch size {
    case "small":
        price *= 0.8
    case "medium":
        price *= 1.0
    case "large":
        price *= 1.3
    }

    price += float64(len(extras)) * 0.5
    return price
}

// POST /orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var order models.Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, "JSON invalide", http.StatusBadRequest)
        return
    }

    var drinkFound *models.Drink
    for _, d := range Drinks {
        if d.ID == order.DrinkID {
            drinkFound = &d
            break
        }
    }

    if drinkFound == nil {
        http.Error(w, "Boisson introuvable", http.StatusBadRequest)
        return
    }

    order.ID = generateOrderID()
    order.DrinkName = drinkFound.Name
    order.Status = models.StatusPending
    order.OrderedAt = time.Now()
    order.TotalPrice = CalculatePrice(drinkFound.BasePrice, order.Size, order.Extras)

    Orders = append(Orders, order)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}

// GET /orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Orders)
}

// GET /orders/{id}
func GetOrder(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.PathValue("id")

    for _, order := range Orders {
        if order.ID == id {
            json.NewEncoder(w).Encode(order)
            return
        }
    }

    http.Error(w, "Commande introuvable", http.StatusNotFound)
}

// PATCH /orders/{id}/status
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    id := r.PathValue("id")

    var payload struct {
        Status string `json:"status"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "JSON invalide", http.StatusBadRequest)
        return
    }

    for i := range Orders {
        if Orders[i].ID == id {
            Orders[i].Status = models.OrderStatus(payload.Status)
            json.NewEncoder(w).Encode(Orders[i])
            return
        }
    }

    http.Error(w, "Commande introuvable", http.StatusNotFound)
}

// DELETE /orders/{id}
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    for i, order := range Orders {
        if order.ID == id {

            if order.Status == models.StatusPickedUp {
                http.Error(w, "Commande déjà récupérée", http.StatusBadRequest)
                return
            }

            Orders = append(Orders[:i], Orders[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Commande introuvable", http.StatusNotFound)
}
