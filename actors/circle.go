package actors

import (
	"image/color"
	"math"

	"kosh/vpaul/floating/components"
	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CircleActor struct {
	core.Actor
	img *ebiten.Image
	op  ebiten.DrawImageOptions
}

func (ca *CircleActor) draw(screen *ebiten.Image, pos utils.Vec) {
	ca.op.GeoM.Reset()
	ca.op.GeoM.Translate(pos.X-core.Radius, pos.Y-core.Radius)
	screen.DrawImage(ca.img, &ca.op)
}

func newCircleImage(clr color.RGBA) *ebiten.Image {
	d := int(math.Ceil(core.Radius * 2))
	img := ebiten.NewImage(d, d)
	r := float32(core.Radius)
	vector.DrawFilledCircle(img, r, r, r, clr, true)
	return img
}

// NewCircle creates a dynamic actor with MovementComponent and CollisionComponent attached.
func NewCircle(x, y float64) *core.Actor {
	ca := &CircleActor{
		img: newCircleImage(color.RGBA{R: 50, G: 200, B: 255, A: 255}),
	}
	ca.Pos = utils.Vec{X: x, Y: y}
	ca.Collider = core.CircleCollider(core.Radius)
	ca.Debug = true
	ca.Draw = ca.draw

	mc := components.NewMovementComponent(&ca.Actor)
	cc := components.NewCollisionComponent(&ca.Actor)
	ca.Components = []core.Updater{mc, cc}
	ca.Input = mc
	return &ca.Actor
}

// NewStatic creates a static collidable circle with no components.
func NewStatic(x, y float64) *core.Actor {
	ca := &CircleActor{
		img: newCircleImage(color.RGBA{R: 50, G: 200, B: 255, A: 255}),
	}
	ca.Pos = utils.Vec{X: x, Y: y}
	ca.Collider = core.CircleCollider(core.Radius)
	ca.Draw = ca.draw
	return &ca.Actor
}
