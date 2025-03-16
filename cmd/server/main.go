package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukegriffith/simulation/internal/sim" // Adjust the import path for your entity package
)

var (
	canvasWidth  float64 = 1000
	canvasHeight float64 = 600
	simMutex     sync.Mutex
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
		simMutex.Lock()
		// Restart the simulation when Enter is pressed
		restartSimulation()
		simMutex.Unlock()
	}
}

func restartSimulation() {
	// Reinitialize the entities or any other necessary state
	sim.InitializeEntities(entityCount, teamCount, canvasWidth, canvasHeight) // You can change the number of entities as needed
	sim.InitializeFood(foodCount, canvasWidth, canvasHeight)
	fmt.Println("Simulation restarted.")
}

var (
	entityCount = 10
	foodCount   = 200
	teamCount   = 2
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
	sim.SetConfig(sim.Config{MinSize: 5, StartMaxSize: 10, MaxSize: 15, BaseSpeed: 10})
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

type MessageType struct {
	Type string
}

type CanvasData struct {
	Width  int
	Height int
}

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
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected:", err)
			delete(clients, ws)
			activeConnections--
			fmt.Println("Active connections:", activeConnections)
			break
		}

		if messageType == websocket.TextMessage {

			var msgType MessageType
			err := json.Unmarshal(message, &msgType)
			if err != nil {
				fmt.Println("unable to marshal canvas data")
			}

			switch t := msgType.Type; t {
			case "resize":
				resize(message)
			case "settings":
				settings(message)
			}
		}
	}
}

type Settings struct {
	Population   int
	TeamCount    int
	FoodCount    int
	MinSize      float64
	StartMaxSize float64
	MaxSize      float64
	BaseSpeed    float64
}

func settings(message []byte) {
	var data Settings
	fmt.Println("settings called")
	err := json.Unmarshal(message, &data)
	if err != nil {
		fmt.Println("unable to marshal settings data", err)
		fmt.Println(string(message))
	}
	teamCount = data.TeamCount
	entityCount = data.Population
	foodCount = data.FoodCount
	simMutex.Lock()
	sim.SetConfig(sim.Config{
		MinSize:      data.MinSize,
		StartMaxSize: data.StartMaxSize,
		MaxSize:      data.MaxSize,
		BaseSpeed:    data.BaseSpeed,
	})
	restartSimulation()
	simMutex.Unlock()

}

func resize(message []byte) {
	var data CanvasData
	fmt.Println("resize called")
	err := json.Unmarshal(message, &data)
	if err != nil {
		fmt.Println("unable to marshal canvas data")
	}
	fmt.Println("set canvas", data)
	canvasWidth = float64(data.Width)
	canvasHeight = float64(data.Height)
	sim.SetCanvas(canvasWidth, canvasHeight)

	simMutex.Lock()
	// Restart the simulation when Enter is pressed
	restartSimulation()
	simMutex.Unlock()
}

func updateSimulationPeriodically() {
	ticker := time.NewTicker(16 * time.Millisecond) // Roughly 60 FPS
	defer ticker.Stop()                             // Ensure the ticker is stopped when the function exits
	previousTime := time.Now()                      // Track the previous time for deltaTime calculation

	for {
		simMutex.Lock()
		<-ticker.C
		currentTime := time.Now()
		deltaTime := currentTime.Sub(previousTime).Seconds() // Calculate deltaTime in seconds
		previousTime = currentTime

		// Skip simulation updates if no active connections
		if activeConnections == 0 {
			simMutex.Unlock()
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
		simMutex.Unlock()
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
