package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lukegriffith/simulation/internal/entity"
)

func main() {
	entity.SetCanvas(800, 500)
	entity.InitializeEntities(5)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/simulate", simulateHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {

	entity.UpdateSimulation()
	// Logic for processing the simulation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.GetEntities())
}
