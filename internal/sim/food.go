package sim

import (
	"fmt"
	"math/rand"
)

type Food struct {
	ID     int     // Unique identifier for the food
	X, Y   float64 // Position of the food
	Size   float64 // Size of the food (affects growth rate)
	Active bool    // Whether the food is still available
}

func InitializeFood(count int, canvasWidth, canvasHeight float64) {
	foods = make([]*Food, count)

	for i := 0; i < count; i++ {
		foods[i] = &Food{
			ID:     i + 1,
			X:      randFloat(0, canvasWidth),
			Y:      randFloat(0, canvasHeight),
			Size:   randFloat(2, 5), // Random size for the food items
			Active: true,
		}
	}
}

func RespawnFood(chance float64, canvasWidth, canvasHeight float64) {
	for i := range foods {
		if !foods[i].Active && rand.Float64() < chance {
			foods[i] = &Food{
				ID:     foods[i].ID,
				X:      randFloat(0, canvasWidth),
				Y:      randFloat(0, canvasHeight),
				Size:   randFloat(2, 5),
				Active: true,
			}
			fmt.Printf("Food %d respawned.\n", foods[i].ID)
		}
	}
}
