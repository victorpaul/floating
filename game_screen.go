package main

import (
	"image/color"
	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/worlds"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  float64 = 800
	ScreenHeight float64 = 600
)

// GameScreen wires the world and controller together and implements ebiten.Game.
type GameScreen struct {
	world      *core.World
	controller *core.Controller
	ctx        *core.GameContext
	debug      *core.DebugOverlay
}

func NewGameScreen() *GameScreen {

	player, world := worlds.NewWorldPool()
	//player, world := worlds.NewWorldSandBox()
	controller := core.NewController(player, ScreenWidth)

	return &GameScreen{
		world:      world,
		controller: controller,
		ctx: &core.GameContext{
			World:      world,
			Controller: controller,
		},
		debug: core.NewDebugOverlay(),
	}
}

func (gs *GameScreen) Update() error {
	start := gs.debug.BeginUpdate()
	gs.controller.Update(gs.ctx)     // apply input → actor velocity
	gs.world.Update(gs.ctx)          // move actors
	gs.controller.LateUpdate(gs.ctx) // snap camera to actor's new position
	gs.debug.EndUpdate(start)
	gs.debug.ActorCount = len(gs.world.Actors)
	return nil
}

func (gs *GameScreen) Draw(screen *ebiten.Image) {
	start := gs.debug.BeginDraw()
	screen.Fill(color.RGBA{R: 10, G: 20, B: 40, A: 255})
	gs.world.Draw(screen, gs.ctx)
	gs.controller.Draw(screen, gs.ctx)
	gs.debug.EndDraw(start)
	gs.debug.Draw(screen)
}

func (gs *GameScreen) Layout(_, _ int) (int, int) {
	return int(ScreenWidth), int(ScreenHeight)
}
