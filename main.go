package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "gotpcoffeeshop/handlers"
    "gotpcoffeeshop/middlewares"
    "gotpcoffeeshop/models"
)

func main() {

    handlers.Drinks = []models.Drink{
        {ID: "DRK-001", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
        {ID: "DRK-002", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
        {ID: "DRK-003", Name: "Latte", Category: "coffee", BasePrice: 4.00},
        {ID: "DRK-004", Name: "Americano", Category: "coffee", BasePrice: 3.00},
        {ID: "DRK-005", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
        {ID: "DRK-006", Name: "Iced Coffee", Category: "coffee", BasePrice: 4.50},
    }

    r := mux.NewRouter()

    r.HandleFunc("/menu", handlers.GetMenu).Methods("GET")
    r.HandleFunc("/menu/{id}", handlers.GetDrink).Methods("GET")

    r.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
    r.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
    r.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")
    r.HandleFunc("/orders/{id}/status", handlers.UpdateOrderStatus).Methods("PATCH")
    r.HandleFunc("/orders/{id}", handlers.DeleteOrder).Methods("DELETE")

    fmt.Println("🚀 API Coffee Shop lancée sur http://localhost:8080")
    http.ListenAndServe(":8080", middlewares.Cors(r))
}
