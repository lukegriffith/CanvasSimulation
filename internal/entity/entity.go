package entity

type Entity struct {
	ID     int     // Unique identifier for the entity
	X, Y   float64 // Position of the entity
	VX, VY float64 // Velocity of the entity
	Width  float64 // Width of the entity (if it's rectangular)
	Height float64 // Height of the entity (if it's rectangular)
	Active bool    // Whether the entity is active in the simulation
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
