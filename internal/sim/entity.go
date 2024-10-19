package sim

import (
	"fmt"
	"math"
)

type State string

const (
	SeekFoodState            State = "SeekFood"
	AssistingTeamMemberState State = "AssistingTeamMember"
	SeekWeakerEnemyState     State = "AssistingWeakerEnemyState"
)

type Entity struct {
	ID                int     // Unique identifier for the entity
	X, Y              float64 // Position of the entity
	VX, VY            float64 // Velocity of the entity
	Width             float64 // Width of the entity (if it's rectangular)
	Height            float64 // Height of the entity (if it's rectangular)
	Active            bool    // Whether the entity is active in the simulation
	Health            float64 // Health of the entity
	MaxHealth         float64
	Invulnerable      bool    // Whether the entity is currently invulnerable
	InvulnTimer       float64 // Time remaining in the invulnerable state
	HungerLevel       float64
	TeamID            int
	TeamNeed          float64
	TeamTimeout       bool
	TeamAssistTimeout float64
	State             State
}

func (e *Entity) DecideAction(entities []*Entity, foods []*Food) {
	// Simple decision criteria
	if e.HungerLevel > 80 {
		// If hunger is critical, prioritize seeking food
		e.SeekFood(foods)
		e.State = SeekFoodState
	} else if e.TeamNeed > 50 && e.TeamAssistTimeout <= 0 {
		// If a teammate needs help, assist the teammate
		e.AssistTeamMember(entities)
		e.State = AssistingTeamMemberState
	} else {
		// Default action
		e.SeekWeakerEnemy(entities)
		e.State = SeekWeakerEnemyState
	}
}

func (e *Entity) SeekWeakerEnemy(entities []*Entity) {
	if !e.Active {
		return // Skip if the entity is not active
	}

	// Find the closest active food
	var closestWeakerEnemy *Entity
	closestDistanceSquared := float64(-1)

	for i := range entities {
		other := entities[i]
		if !other.Active || other.ID == e.ID || other.TeamID == e.TeamID || other.Width >= e.Width {
			continue // Skip inactive food
		}

		// Calculate squared distance to the e
		dx := e.X - e.X
		dy := e.Y - e.Y
		distanceSquared := dx*dx + dy*dy

		// Check if this is the closest e found so far
		if closestDistanceSquared < 0 || distanceSquared < closestDistanceSquared {
			closestDistanceSquared = distanceSquared
			closestWeakerEnemy = other
		}
	}

	// Move towards the closest e if found
	if closestWeakerEnemy != nil {
		dx := closestWeakerEnemy.X - e.X
		dy := closestWeakerEnemy.Y - e.Y

		// Normalize the direction and adjust velocity using pythagoras
		length := math.Sqrt(dx*dx + dy*dy)
		if length != 0 {
			// TODO: Scale on hunger need?
			e.VX += (dx / length) * 0.1 // Scale speed towards the food
			e.VY += (dy / length) * 0.1
		}
	}
}

func (e *Entity) EvaluateTeamNeed(entities []*Entity) {
	// Reset TeamNeed to avoid carrying over previous values
	e.TeamNeed = 0

	// Iterate over all entities to evaluate teammates' needs
	for _, teammate := range entities {
		// Check if the entity is on the same team, is not itself, and has low health
		if teammate.TeamID == e.TeamID && teammate != e && teammate.Health < 50 {
			// Calculate the distance to the teammate
			distance := e.DistanceTo(teammate)

			// If the teammate is within a certain range, increase the TeamNeed score
			if distance < 100.0 { // Example range threshold
				// Increase TeamNeed based on the severity of the teammate's condition
				e.TeamNeed += 50.0 - teammate.Health // The lower the health, the higher the need
			}
		}
	}

	// Cap TeamNeed to a maximum value if needed
	if e.TeamNeed > 100 {
		e.TeamNeed = 100
	}
}

func (e *Entity) Sense(nearbyEntities []Entity, canvasWidth, canvasHeight float64) {
	for _, other := range nearbyEntities {
		if other.ID == e.ID || other.TeamID == e.TeamID || !other.Active {
			continue // Skip self or same team or inactive entities
		}

		// Calculate the distance to the other entity
		dx := other.X - e.X
		dy := other.Y - e.Y
		distance := dx*dx + dy*dy // Squared distance for efficiency

		// Flee behavior: if the other entity is larger, move away
		if other.Width > e.Width {
			if distance < 10000 { // Threshold for fleeing behavior
				e.VX -= dx / 50 // Move away from the other entity
				e.VY -= dy / 50
			}
		} else {
			// Seek behavior: if the other entity is smaller, move towards it
			if distance < 40000 { // Threshold for seeking behavior
				e.VX += dx / 50 // Move towards the other entity
				e.VY += dy / 50
			}
		}
	}

	// Sense the boundaries and adjust velocity to avoid going out of bounds
	if e.X < 50 {
		e.VX += 1 // Move right
	} else if e.X+e.Width > canvasWidth-50 {
		e.VX -= 1 // Move left
	}
	if e.Y < 50 {
		e.VY += 1 // Move down
	} else if e.Y+e.Height > canvasHeight-50 {
		e.VY -= 1 // Move up
	}
}

func (e *Entity) Act(nearbyEntities []*Entity, canvasWidth, canvasHeight, deltaTime float64) {
	// Step 1: Check if the entity is active
	if !e.Active {
		return
	}

	// Step 2: Handle the invulnerability timer
	if e.Invulnerable {
		e.InvulnTimer -= deltaTime
		if e.InvulnTimer <= 0 {
			e.Invulnerable = false
			e.InvulnTimer = 0
			fmt.Printf("Entity %d is no longer invulnerable.\n", e.ID)
		}
		return
	}

	if e.TeamTimeout {
		e.TeamAssistTimeout -= deltaTime
		if e.TeamAssistTimeout <= 0 {
			e.TeamTimeout = false
			e.TeamAssistTimeout = 0
			fmt.Printf("Entity %d is no longer on team timeout.\n", e.ID)
		}
	}

	// Step 3: Update position based on velocity
	e.X += e.VX * deltaTime
	e.Y += e.VY * deltaTime

	// Step 4: Limit the speed based on the size of the entity
	baseSpeed := 10.0
	sizeFactor := 1.0 / (1.0 + (e.Width / 100.0)) // Speed decreases as size increases
	maxSpeed := baseSpeed * sizeFactor

	// Cap the velocity components to the maximum speed
	e.VX = clamp(e.VX, -maxSpeed, maxSpeed)
	e.VY = clamp(e.VY, -maxSpeed, maxSpeed)

	// Step 5: Keep the entity within the canvas boundaries
	if e.X < 0 {
		e.X = 0
		e.VX = -e.VX // Reverse direction upon hitting the left boundary
	} else if e.X+e.Width > canvasWidth {
		e.X = canvasWidth - e.Width
		e.VX = -e.VX // Reverse direction upon hitting the right boundary
	}

	if e.Y < 0 {
		e.Y = 0
		e.VY = -e.VY // Reverse direction upon hitting the top boundary
	} else if e.Y+e.Height > canvasHeight {
		e.Y = canvasHeight - e.Height
		e.VY = -e.VY // Reverse direction upon hitting the bottom boundary
	}

	// Step 6: Interact with nearby entities (consume behavior)
	e.Consume(nearbyEntities)

	if e.HungerLevel > 0 {
		e.HungerLevel -= 1.0
	}
	// Step 7: Deactivate if health is depleted
	if e.Health <= 0 {
		e.Active = false
		fmt.Printf("Entity %d has been deactivated due to depleted health.\n", e.ID)
	}

}

// Helper function to clamp a value between a minimum and maximum range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (e *Entity) Consume(nearbyEntities []*Entity) {
	for i := range nearbyEntities {
		other := nearbyEntities[i]

		// Skip self, inactive entities, or if the other entity is invulnerable or same team
		if other.ID == e.ID || !other.Active || other.Invulnerable || other.TeamID == e.TeamID {
			continue
		}

		// Step 1: Check if the other entity is smaller
		if other.Width >= e.Width {
			continue // If the other entity is not smaller, skip it
		}

		// Step 2: Check if the other entity is close enough to be consumed
		dx := other.X - e.X
		dy := other.Y - e.Y
		distanceSquared := dx*dx + dy*dy
		entityRadius := e.Width
		consumptionRange := 1.5 * entityRadius
		consumptionThreshold := consumptionRange * consumptionRange

		if distanceSquared > consumptionThreshold {
			continue // If the other entity is too far, skip it
		}

		// Step 3: Consume the other entity
		// Apply a health penalty scaled by the relative size of the entities
		healthPenalty := 40.0 * (other.Width / e.Width)
		// Heal entity for eating
		// Damage other for being consumed
		e.Health -= healthPenalty * 0.3
		other.Health -= healthPenalty

		// If health drops below zero, deactivate the entity
		if e.Health <= 0 {
			e.Active = false
			fmt.Printf("Entity %d has died after consuming Entity %d.\n", e.ID, other.ID)
			return // Stop processing further consumption for this entity
		}

		// Increase the size of the consuming entity
		e.Grow(0.1)

		// Step 4: Trigger invulnerability
		other.Invulnerable = true
		other.InvulnTimer = 1.0 + (e.Width / 200.0) // Set invulnerability duration based on size

		// Logging for debugging
		fmt.Printf("Entity %d (Team %d) consumed Entity %d (Team %d). Healing %f, New size: (%.2f, %.2f). Health: %.2f\n", e.ID, e.TeamID, other.ID, other.TeamID, healthPenalty, e.Width, e.Height, e.Health)
		fmt.Printf("Entity %d (Team %d) consumed by Entity %d (Team %d). Taking %.2f damage. Health: %.2f\n", other.ID, other.TeamID, e.ID, e.TeamID, healthPenalty, e.Health)
	}
}

func (e *Entity) SeekFood(nearbyFood []*Food) {
	if !e.Active {
		return // Skip if the entity is not active
	}

	// Find the closest active food
	var closestFood *Food
	closestDistanceSquared := float64(-1)

	for i := range nearbyFood {
		food := nearbyFood[i]
		if !food.Active {
			continue // Skip inactive food
		}

		// Calculate squared distance to the food
		dx := food.X - e.X
		dy := food.Y - e.Y
		distanceSquared := dx*dx + dy*dy

		// Check if this is the closest food found so far
		if closestDistanceSquared < 0 || distanceSquared < closestDistanceSquared {
			closestDistanceSquared = distanceSquared
			closestFood = food
		}
	}

	// Move towards the closest food if found
	if closestFood != nil {
		dx := closestFood.X - e.X
		dy := closestFood.Y - e.Y

		// Normalize the direction and adjust velocity using pythagoras
		length := math.Sqrt(dx*dx + dy*dy)
		if length != 0 {
			// TODO: Scale on hunger need?
			e.VX += (dx / length) * 0.1 // Scale speed towards the food
			e.VY += (dy / length) * 0.1
		}
	}
}

func (e *Entity) ConsumeFood(nearbyFood []*Food) {
	if !e.Active {
		return // Skip if the entity is not active
	}

	for i := range nearbyFood {
		food := nearbyFood[i]
		if !food.Active {
			continue // Skip inactive food
		}

		// Check if the entity is close enough to consume the food
		dx := food.X - e.X
		dy := food.Y - e.Y
		distanceSquared := dx*dx + dy*dy
		foodThreshold := (food.Size + e.Width) * 1.2 // Consumption range based on entity size

		if distanceSquared < foodThreshold*foodThreshold {
			// Consume the food
			e.Grow(0.1)
			e.Health += food.Size * 2
			food.Active = false // Deactivate the food

			fmt.Printf("Entity %d consumed Food %d and grew.\n", e.ID, food.ID)
			e.HungerLevel = 100.0
			break // Only consume one food per update
		}
	}
}
func (e *Entity) SetActive(active bool) {
	e.Active = active
}

func (e *Entity) Grow(factor float64) {
	if e.Width < config.MaxSize {
		e.Width += e.Width * factor
	}
}

func (e *Entity) AssistTeamMember(entities []*Entity) {
	// Find the nearest teammate in need of help
	var nearestTeammate *Entity
	minDistance := 4.0

	for _, teammate := range entities {
		if teammate.TeamID == e.TeamID && teammate != e && teammate.Health < 50 {
			// Calculate the distance to the teammate
			distance := e.DistanceTo(teammate)

			// Find the closest teammate in need of assistance
			if distance < minDistance {
				minDistance = distance
				nearestTeammate = teammate
			}
		}
	}

	// If a teammate is found, perform an assist action
	if nearestTeammate != nil {
		e.PerformAssistAction(nearestTeammate)
		e.TeamTimeout = true
		e.TeamAssistTimeout = 5.0
	}
}

// Calculate distance between two entities (helper method)
func (e *Entity) DistanceTo(other *Entity) float64 {
	dx := e.X - other.X
	dy := e.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Perform an assist action, e.g., healing the teammate (helper method)
func (e *Entity) PerformAssistAction(teammate *Entity) {
	// Example: Heal the teammate by a certain amount
	healAmount := 0.1
	teammate.Health += healAmount

	// Cap the health at a maximum value (assuming 100 is the max)
	if teammate.Health > 100 {
		teammate.Health = 100
	}

	// Optionally, you could also print a log or message
	fmt.Printf("Entity %d assisted teammate %d, healing them by %f.\n", e.ID, teammate.ID, healAmount)
}
