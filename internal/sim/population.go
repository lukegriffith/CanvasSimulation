package sim

import (
	"math/rand"
	"time"
)

var entities []*Entity
var foods []*Food

var config Config

type Config struct {
	MinSize, StartMaxSize, MaxSize float64
}

func InitializeEntities(population int, teams int, canvasWidth float64, canvasHeight float64) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	entities = make([]*Entity, population) // Create a slice to hold the entities
	var teamCounter = 0
	for i := 0; i < population; i++ {
		entities[i] = &Entity{
			ID:        i + 1,
			X:         randFloat(0, canvasWidth),                      // Random X position between 0 and 800
			Y:         randFloat(0, canvasHeight),                     // Random Y position between 0 and 600
			VX:        randFloat(-10, 10),                             // Random velocity X between -2 and 2
			VY:        randFloat(-10, 10),                             // Random velocity Y between -2 and 2
			Width:     randFloat(config.MinSize, config.StartMaxSize), // Random width between 20 and 100
			Active:    true,
			Health:    100, // Set initial health to 100
			MaxHealth: 100,
			TeamID:    teamCounter % teams,
		}
		teamCounter = teamCounter + 1
	}
}

// Helper function to generate a random float64 between min and max
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

var respawnTimer float64

func UpdateSimulation(deltaTime float64) {
	for i := range entities {
		if entities[i].Active {
			// Evaluate team needs to update the entity's priority
			entities[i].EvaluateTeamNeed(entities)

			// Decide on the action (assist teammate, seek food, etc.)
			entities[i].DecideAction(entities, foods)

			// Consume food if possible
			entities[i].ConsumeFood(foods)

			// Update position, perform other actions, and keep within the canvas
			entities[i].Act(entities, canvasWidth, canvasHeight, deltaTime)
		}
	}

	// Periodically respawn food items with a certain chance
	RespawnFood(0.001, canvasWidth, canvasHeight)

	if respawnTimer >= 5.0 {
		RespawnFood(0.01, canvasWidth, canvasHeight)
		respawnTimer = 0.0
	}
}

func GetEntities() []*Entity {
	return entities
}

func GetFood() []*Food {
	return foods
}

var (
	canvasWidth, canvasHeight float64
)

func SetCanvas(w float64, h float64) {
	canvasWidth = w
	canvasHeight = h
}

func SetConfig(c Config) {
	config = c
}
