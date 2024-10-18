package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukegriffith/simulation/internal/sim" // Adjust the import path for your entity package
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []sim.Entity)      // Broadcast channel for entities

func main() {
	sim.SetCanvas(800, 500)
	sim.InitializeEntities(100)

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the simulation update loop in a separate goroutine
	go updateSimulationPeriodically()
	go handleMessages()

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	// Clean up when the client disconnects
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected:", err)
			delete(clients, ws)
			break
		}
	}
}

func updateSimulationPeriodically() {
	ticker := time.NewTicker(16 * time.Millisecond) // Roughly 60 FPS
	for {
		<-ticker.C
		sim.UpdateSimulation()        // Update the simulation
		entities := sim.GetEntities() // Get the current state of entities

		fmt.Println("Broadcasting entities:", entities) // Debug log

		broadcast <- entities // Send the updated entities to the broadcast channel
	}
}

func handleMessages() {
	for {
		entities := <-broadcast
		for client := range clients {
			err := client.WriteJSON(entities)
			if err != nil {
				fmt.Println("Error sending data to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
