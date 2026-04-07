package actors

import (
	"image/color"

	"kosh/vpaul/floating/components"
	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func circleDraw(clr color.RGBA) core.DrawFunc {
	return func(screen *ebiten.Image, pos utils.Vec) {
		vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y),
			float32(core.Radius), clr, true)
	}
}

// New creates a dynamic actor with MovementComponent and CollisionComponent attached.
func NewCircle(x, y float64) *core.Actor {
	a := &core.Actor{
		Pos:      utils.Vec{X: x, Y: y},
		Collider: core.CircleCollider(core.Radius),
		Debug:    true,
		Draw:     circleDraw(color.RGBA{R: 50, G: 200, B: 255, A: 255}),
	}
	mc := components.NewMovementComponent(a)
	cc := components.NewCollisionComponent(a)

	a.Components = []core.Updater{mc, cc}
	a.Input = mc
	return a
}

// NewStatic creates a static collidable circle with no components.
func NewStatic(x, y float64) *core.Actor {
	return &core.Actor{
		Pos:      utils.Vec{X: x, Y: y},
		Collider: core.CircleCollider(core.Radius),
		Draw:     circleDraw(color.RGBA{R: 50, G: 200, B: 255, A: 255}),
	}
}
