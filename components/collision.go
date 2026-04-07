package components

import (
	"math"

	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"
)

// CollisionComponent pushes the owner actor away from any overlapping colliders
// in the world. The separation force magnitude equals the owner's current speed,
// producing smooth sliding rather than hard stopping.
// The owner is assumed to have a circle collider.
type CollisionComponent struct {
	core.Component
	Owner *core.Actor
}

func NewCollisionComponent(owner *core.Actor) *CollisionComponent {
	return &CollisionComponent{Owner: owner}
}

func (c *CollisionComponent) Update(ctx *core.GameContext) {
	owner := c.Owner
	if owner.Collider == nil || owner.Collider.Shape != core.ShapeCircle {
		return
	}

	for _, other := range ctx.World.Grid.GetNearby(owner.Pos) {
		if other == owner || other.Collider == nil {
			continue
		}

		pushDir, overlap, ok := resolveOverlap(owner.Pos, owner.Collider.Radius, other)
		if !ok {
			continue
		}

		// Position correction: push owner fully out of overlap this frame
		owner.Pos = owner.Pos.Add(pushDir.Scale(overlap))

		// If pushed upward, the actor is standing on something
		if pushDir.Y < -0.3 {
			owner.Grounded = true
		}

		// Velocity projection: cancel only the component driving into the surface
		relV := owner.Vel.Dot(pushDir)
		if relV < 0 {
			owner.Vel = owner.Vel.Sub(pushDir.Scale(relV))
		}
	}
}

// resolveOverlap returns the push direction, overlap depth, and whether an overlap
// exists between a circle (pos, radius) and another actor's collider.
func resolveOverlap(pos utils.Vec, radius float64, other *core.Actor) (pushDir utils.Vec, overlap float64, ok bool) {
	switch other.Collider.Shape {

	case core.ShapeCircle:
		diff := pos.Sub(other.Pos)
		dist := diff.Length()
		minDist := radius + other.Collider.Radius
		if dist == 0 || dist >= minDist {
			return
		}
		return diff.Normalized(), minDist - dist, true

	case core.ShapeBox:
		half := other.Collider.HalfSize
		// Transform circle center into box local space (inverse rotation)
		local := pos.Sub(other.Pos).Rotate(-other.Collider.Angle)
		// Closest point on box to circle center
		cx := clamp(local.X, -half.X, half.X)
		cy := clamp(local.Y, -half.Y, half.Y)
		localDiff := utils.Vec{X: local.X - cx, Y: local.Y - cy}
		dist := localDiff.Length()

		if dist >= radius {
			return
		}

		var localPush utils.Vec
		if dist > 0 {
			localPush = localDiff.Normalized()
		} else {
			// Circle center is inside the box — push out via nearest edge
			dx := half.X - math.Abs(local.X)
			dy := half.Y - math.Abs(local.Y)
			if dx < dy {
				localPush = utils.Vec{X: math.Copysign(1, local.X)}
			} else {
				localPush = utils.Vec{Y: math.Copysign(1, local.Y)}
			}
		}

		// Rotate push direction back to world space
		return localPush.Rotate(other.Collider.Angle), radius - dist, true
	}
	return
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
