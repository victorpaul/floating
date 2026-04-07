package actors

import (
	"image/color"
	"math"

	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func boxDraw(halfW, halfH, angle float64, clr color.RGBA) core.DrawFunc {
	half := utils.Vec{X: halfW, Y: halfH}
	return func(screen *ebiten.Image, pos utils.Vec) {
		core.DrawRotatedFilledRect(screen, float32(pos.X), float32(pos.Y), half, angle, clr)
	}
}

// NewStaticSquare creates a static collidable box with no components.
// halfW and halfH are the half-extents of the box.
func NewStaticSquare(x, y, halfW, halfH float64) *core.Actor {
	return &core.Actor{
		Pos:      utils.Vec{X: x, Y: y},
		Collider: core.BoxCollider(halfW, halfH),
		Draw:     boxDraw(halfW, halfH, 0, color.RGBA{R: 200, G: 150, B: 50, A: 255}),
	}
}

// NewStaticRotatedSquare creates a static collidable rotated box with no components.
// angle is in degrees.
func NewStaticRotatedSquare(x, y, halfW, halfH, degrees float64) *core.Actor {
	angle := degrees * math.Pi / 180
	return &core.Actor{
		Pos:      utils.Vec{X: x, Y: y},
		Collider: core.RotatedBoxCollider(halfW, halfH, angle),
		Draw:     boxDraw(halfW, halfH, angle, color.RGBA{R: 200, G: 150, B: 50, A: 255}),
	}
}
