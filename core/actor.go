package core

import (
	"image/color"

	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawFunc is called during Draw with the screen and the actor's screen-space position.
type DrawFunc func(screen *ebiten.Image, screenPos utils.Vec)

const Radius float64 = 18

// Actor is a game entity with position, velocity, and an optional set of components.
type Actor struct {
	Pos        utils.Vec
	Vel        utils.Vec
	Collider   *Collider     // nil = not collidable
	Input      InputReceiver // nil for non-player actors
	Components []Updater
	Grounded   bool     // set by CollisionComponent, read by MovementComponent for jump gating
	Debug      bool
	Draw       DrawFunc // set by actor constructors; called each frame with screen-space position
}

func (a *Actor) Update(ctx *GameContext) {
	for _, c := range a.Components {
		c.Update(ctx)
	}
}

func (a *Actor) Render(screen *ebiten.Image, ctx *GameContext) {
	s := a.Pos.Sub(ctx.Camera)

	if a.Draw != nil {
		a.Draw(screen, s)
	}

	if a.Debug && a.Collider != nil {
		switch a.Collider.Shape {
		case ShapeCircle:
			vector.StrokeCircle(screen, float32(s.X), float32(s.Y),
				float32(a.Collider.Radius), 1,
				color.RGBA{R: 255, A: 255}, true)
		case ShapeBox:
			StrokeRotatedRect(screen, float32(s.X), float32(s.Y),
				a.Collider.HalfSize, a.Collider.Angle, 1,
				color.RGBA{R: 255, A: 255})
		}
	}
}
