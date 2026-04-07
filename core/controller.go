package core

import (
	"os"

	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Controller reads player input and drives a single controlled actor.
// It also owns the camera position, written into GameContext each tick.
type Controller struct {
	Camera      utils.Vec
	Controlled  *Actor
	screenWidth float64
}

func NewController(target *Actor, screenWidth float64) *Controller {
	return &Controller{
		Controlled:  target,
		screenWidth: screenWidth,
	}
}

// Update reads input and applies it to the controlled actor's velocity.
func (c *Controller) Update(ctx *GameContext) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		ctx.World.DebugGrid = !ctx.World.DebugGrid
		for _, a := range ctx.World.Actors {
			a.Debug = ctx.World.DebugGrid
			for _, comp := range a.Components {
				if d, ok := comp.(Debuggable); ok {
					d.SetDebug(a.Debug)
				}
			}
		}
	}

	actor := c.Controlled
	if actor == nil {
		return
	}

	if actor.Input != nil {
		actor.Input.ApplyInput(MovementInput{
			Left:  ebiten.IsKeyPressed(ebiten.KeyA),
			Right: ebiten.IsKeyPressed(ebiten.KeyD),
			Up:    inpututil.IsKeyJustPressed(ebiten.KeyW),
		})
	}
}

// LateUpdate snaps the camera to the controlled actor's current (already moved) position.
// Call after world.Update each frame.
func (c *Controller) LateUpdate(ctx *GameContext) {
	actor := c.Controlled
	if actor == nil {
		return
	}

	c.Camera.X = actor.Pos.X - c.screenWidth/2
	c.Camera.Y = 0

	ctx.Camera = c.Camera
}

// Draw renders controller-owned UI elements (currently none).
func (c *Controller) Draw(screen *ebiten.Image, ctx *GameContext) {}
