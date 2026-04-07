package actors

import (
	"image/color"
	"math"

	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type SquareActor struct {
	core.Actor
	img      *ebiten.Image
	op       ebiten.DrawImageOptions
	baseGeoM ebiten.GeoM
}

func (sa *SquareActor) draw(screen *ebiten.Image, pos utils.Vec) {
	sa.op.GeoM = sa.baseGeoM
	sa.op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(sa.img, &sa.op)
}

func newSquareActor(x, y, halfW, halfH, angle float64, collider *core.Collider) *SquareActor {
	w := int(math.Ceil(halfW * 2))
	h := int(math.Ceil(halfH * 2))
	img := ebiten.NewImage(w, h)
	img.Fill(color.RGBA{R: 200, G: 150, B: 50, A: 255})

	sa := &SquareActor{img: img}
	sa.baseGeoM.Translate(-halfW, -halfH)
	sa.baseGeoM.Rotate(angle)

	sa.Pos = utils.Vec{X: x, Y: y}
	sa.Collider = collider
	sa.Draw = sa.draw
	return sa
}

// NewStaticSquare creates a static collidable box with no components.
func NewStaticSquare(x, y, halfW, halfH float64) *core.Actor {
	sa := newSquareActor(x, y, halfW, halfH, 0, core.BoxCollider(halfW, halfH))
	return &sa.Actor
}

// NewStaticRotatedSquare creates a static collidable rotated box with no components.
// angle is in degrees.
func NewStaticRotatedSquare(x, y, halfW, halfH, degrees float64) *core.Actor {
	angle := degrees * math.Pi / 180
	sa := newSquareActor(x, y, halfW, halfH, angle, core.RotatedBoxCollider(halfW, halfH, angle))
	return &sa.Actor
}
