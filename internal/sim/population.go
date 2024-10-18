package sim

import (
	"math/rand"
	"time"
)

var entities []Entity

func InitializeEntities(population int) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	entities = make([]Entity, population) // Create a slice to hold the entities

	for i := 0; i < population; i++ {
		entities[i] = Entity{
			ID:     i + 1,
			X:      randFloat(0, 800), // Random X position between 0 and 800
			Y:      randFloat(0, 600), // Random Y position between 0 and 600
			VX:     randFloat(-2, 2),  // Random velocity X between -2 and 2
			VY:     randFloat(-2, 2),  // Random velocity Y between -2 and 2
			Width:  randFloat(3, 20),  // Random width between 20 and 100
			Height: randFloat(2, 20),  // Random height between 20 and 100
			Active: true,
		}
	}
}

// Helper function to generate a random float64 between min and max
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func UpdateSimulation() {
	for i := range entities {
		if entities[i].Active {
			// Sense the environment and adjust velocity
			entities[i].Sense(entities, canvasWidth, canvasHeight)
			// Update the position and apply actions
			entities[i].Act()
		}
	}
}

func GetEntities() []Entity {
	return entities
}

var (
	canvasWidth, canvasHeight float64
)

func SetCanvas(w float64, h float64) {
	canvasWidth = w
	canvasHeight = h
}
