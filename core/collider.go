package core

import "kosh/vpaul/floating/utils"

type ShapeType int

const (
	ShapeCircle ShapeType = iota
	ShapeBox
)

// Collider describes the collision volume for an actor.
// Actors with a nil Collider are ignored by collision checks.
type Collider struct {
	Shape    ShapeType
	Radius   float64   // used when Shape == ShapeCircle
	HalfSize utils.Vec // used when Shape == ShapeBox (half-extents)
	Angle    float64   // rotation in radians, used when Shape == ShapeBox
}

func CircleCollider(radius float64) *Collider {
	return &Collider{Shape: ShapeCircle, Radius: radius}
}

func BoxCollider(halfW, halfH float64) *Collider {
	return &Collider{Shape: ShapeBox, HalfSize: utils.Vec{X: halfW, Y: halfH}}
}

func RotatedBoxCollider(halfW, halfH, angle float64) *Collider {
	return &Collider{Shape: ShapeBox, HalfSize: utils.Vec{X: halfW, Y: halfH}, Angle: angle}
}
