package core

// Updater is implemented by any component that needs to run logic each game tick.
type Updater interface {
	Update(*GameContext)
}

// InputReceiver is implemented by components that accept player directional input.
type InputReceiver interface {
	ApplyInput(MovementInput)
}

// MovementInput holds the raw directional intent for a single frame.
type MovementInput struct {
	Left, Right, Up, Down bool
}

// Debuggable is implemented by types that expose a debug toggle.
type Debuggable interface {
	SetDebug(bool)
}

// Component is the base struct embedded by all concrete components.
// Provides shared state and a default no-op Update.
type Component struct {
	Debug bool
}

func (c *Component) Update(_ *GameContext) {}

func (c *Component) SetDebug(d bool) { c.Debug = d }
