package sim

type Entity struct {
	ID     int     // Unique identifier for the entity
	X, Y   float64 // Position of the entity
	VX, VY float64 // Velocity of the entity
	Width  float64 // Width of the entity (if it's rectangular)
	Height float64 // Height of the entity (if it's rectangular)
	Active bool    // Whether the entity is active in the simulation
}

// Sense processes sensory information and updates the entity's internal state
func (e *Entity) Sense(nearbyEntities []Entity, canvasWidth, canvasHeight float64) {
	// Example: Sense the distance to other entities and the boundaries
	for _, other := range nearbyEntities {
		if other.ID == e.ID || !other.Active {
			continue // Skip self or inactive entities
		}

		// Calculate the distance to the other entity
		dx := other.X - e.X
		dy := other.Y - e.Y
		distance := (dx*dx + dy*dy) // Squared distance for simplicity

		// Example sensory processing: if too close to another entity, adjust velocity
		if distance < 10000 { // Threshold for "too close" (squared distance)
			e.VX -= dx / 50 // Move away from the other entity
			e.VY -= dy / 50
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

// Act updates the entity's velocity based on the results of Sense()
func (e *Entity) Act() {
	// Example action: limit the speed to prevent excessive velocity
	maxSpeed := 5.0
	if e.VX > maxSpeed {
		e.VX = maxSpeed
	} else if e.VX < -maxSpeed {
		e.VX = -maxSpeed
	}
	if e.VY > maxSpeed {
		e.VY = maxSpeed
	} else if e.VY < -maxSpeed {
		e.VY = -maxSpeed
	}

	// Update the entity's position based on the velocity
	e.X += e.VX
	e.Y += e.VY
}

func (e *Entity) UpdatePosition() {
	// Update position based on velocity
	e.X += e.VX
	e.Y += e.VY

	// Check for collision with the left or right edge
	if e.X < 0 {
		e.X = 0
		e.VX = -e.VX // Reverse velocity in the X direction
	} else if e.X+e.Width > canvasWidth {
		e.X = canvasWidth - e.Width
		e.VX = -e.VX // Reverse velocity in the X direction
	}

	// Check for collision with the top or bottom edge
	if e.Y < 0 {
		e.Y = 0
		e.VY = -e.VY // Reverse velocity in the Y direction
	} else if e.Y+e.Height > canvasHeight {
		e.Y = canvasHeight - e.Height
		e.VY = -e.VY // Reverse velocity in the Y direction
	}
}

func (e *Entity) CollidesWith(other *Entity) bool {
	return e.X < other.X+other.Width &&
		e.X+e.Width > other.X &&
		e.Y < other.Y+other.Height &&
		e.Y+e.Height > other.Y
}

func (e *Entity) SetActive(active bool) {
	e.Active = active
}
