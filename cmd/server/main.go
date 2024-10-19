package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukegriffith/simulation/internal/sim" // Adjust the import path for your entity package
)

func listenForEnter() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Press Enter to restart the simulation...\n")
		_, err := reader.ReadString('\n') // Wait for Enter key
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// Restart the simulation when Enter is pressed
		restartSimulation()
	}
}

func restartSimulation() {
	// Reinitialize the entities or any other necessary state
	sim.InitializeEntities(entityCount, teamCount, canvasWidth, canvasHeight) // You can change the number of entities as needed
	fmt.Println("Simulation restarted.")
}

const (
	entityCount  = 25
	foodCount    = 200
	teamCount    = 2
	canvasWidth  = 1000
	canvasHeight = 600
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan responseData)      // Broadcast channel for entities

func main() {
	sim.SetCanvas(canvasWidth, canvasHeight)
	sim.SetConfig(sim.Config{5, 10, 25})
	sim.InitializeEntities(entityCount, teamCount, canvasWidth, canvasHeight)
	sim.InitializeFood(foodCount, canvasWidth, canvasHeight)

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the simulation update loop in a separate goroutine
	go listenForEnter()
	go updateSimulationPeriodically()
	go handleMessages()

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

var activeConnections int // Track the number of active connections

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true
	activeConnections++
	fmt.Println("New WebSocket connection established. Active connections:", activeConnections)

	// Clean up when the client disconnects
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected:", err)
			delete(clients, ws)
			activeConnections--
			fmt.Println("Active connections:", activeConnections)
			break
		}
	}
}

func updateSimulationPeriodically() {
	ticker := time.NewTicker(16 * time.Millisecond) // Roughly 60 FPS
	defer ticker.Stop()                             // Ensure the ticker is stopped when the function exits
	previousTime := time.Now()                      // Track the previous time for deltaTime calculation

	for {
		<-ticker.C
		currentTime := time.Now()
		deltaTime := currentTime.Sub(previousTime).Seconds() // Calculate deltaTime in seconds
		previousTime = currentTime

		// Skip simulation updates if no active connections
		if activeConnections == 0 {
			continue
		}

		// Only update the simulation if there are active connections
		if activeConnections > 0 {
			sim.UpdateSimulation(deltaTime) // Pass deltaTime to the UpdateSimulation function
			entities := sim.GetEntities()   // Get the current state of entities
			foods := sim.GetFood()
			// Broadcast the updated entities
			broadcast <- responseData{
				Entities:  entities,
				Foods:     foods,
				TeamCount: teamCount,
			}
		}
	}
}

type responseData struct {
	Entities  []*sim.Entity
	Foods     []*sim.Food
	TeamCount int
}

func handleMessages() {
	for {
		data := <-broadcast
		for client := range clients {
			err := client.WriteJSON(data)
			if err != nil {
				fmt.Println("Error sending data to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
