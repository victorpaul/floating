package main

import (
	"image/color"
	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/water"
	"kosh/vpaul/floating/worlds"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	waterGrid  *water.Grid
}

func NewGameScreen() *GameScreen {
	player, world := worlds.NewWorldPool()
	//player, world := worlds.NewWorldSandBox()
	controller := core.NewController(player, ScreenWidth)

	wg := water.NewGrid(2000, 800, 0, 0)
	wg.BakeGeometry(world.Actors)

	return &GameScreen{
		world:      world,
		controller: controller,
		ctx: &core.GameContext{
			World:      world,
			Controller: controller,
		},
		debug:     core.NewDebugOverlay(),
		waterGrid: wg,
	}
}

func (gs *GameScreen) Update() error {
	start := gs.debug.BeginUpdate()
	gs.controller.Update(gs.ctx)     // apply input → actor velocity
	gs.world.Update(gs.ctx)          // move actors
	gs.controller.LateUpdate(gs.ctx) // snap camera to actor's new position
	gs.waterGrid.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		p := gs.controller.Controlled.Pos
		gs.waterGrid.AddWater(p.X, p.Y, 5)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		p := gs.controller.Controlled.Pos
		gs.waterGrid.PunchHole(p.X, p.Y-40, 3)
	}

	gs.debug.EndUpdate(start)
	gs.debug.ActorCount = len(gs.world.Actors)
	return nil
}

func (gs *GameScreen) Draw(screen *ebiten.Image) {
	start := gs.debug.BeginDraw()
	screen.Fill(color.RGBA{R: 10, G: 20, B: 40, A: 255})
	gs.world.Draw(screen, gs.ctx)
	gs.waterGrid.Draw(screen, gs.ctx.Camera)
	gs.controller.Draw(screen, gs.ctx)
	gs.debug.EndDraw(start)
	gs.debug.Draw(screen)
}

func (gs *GameScreen) Layout(_, _ int) (int, int) {
	return int(ScreenWidth), int(ScreenHeight)
}
