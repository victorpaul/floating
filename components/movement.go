package components

import "kosh/vpaul/floating/core"

const (
	defaultSpeed        = 3.0
	defaultGravity      = 0.4
	defaultJumpForce    = -9.0
	defaultMaxFallSpeed = 12.0
)

// MovementComponent handles horizontal movement, gravity, and jumping.
type MovementComponent struct {
	core.Component
	Owner        *core.Actor
	Speed        float64
	Gravity      float64
	JumpForce    float64
	MaxFallSpeed float64
}

func NewMovementComponent(owner *core.Actor) *MovementComponent {
	return &MovementComponent{
		Owner:        owner,
		Speed:        defaultSpeed,
		Gravity:      defaultGravity,
		JumpForce:    defaultJumpForce,
		MaxFallSpeed: defaultMaxFallSpeed,
	}
}

// Update resets Grounded, applies gravity, then moves the actor.
// Grounded is reset here so CollisionComponent (which runs after) can set it.
func (m *MovementComponent) Update(_ *core.GameContext) {
	m.Owner.Grounded = false

	m.Owner.Vel.Y += m.Gravity
	if m.Owner.Vel.Y > m.MaxFallSpeed {
		m.Owner.Vel.Y = m.MaxFallSpeed
	}

	m.Owner.Pos = m.Owner.Pos.Add(m.Owner.Vel)
}

// ApplyInput sets horizontal velocity and triggers a jump when grounded.
func (m *MovementComponent) ApplyInput(input core.MovementInput) {
	m.Owner.Vel.X = 0
	if input.Left {
		m.Owner.Vel.X = -m.Speed
	}
	if input.Right {
		m.Owner.Vel.X = m.Speed
	}
	if input.Up && m.Owner.Grounded {
		m.Owner.Vel.Y = m.JumpForce
	}
}
