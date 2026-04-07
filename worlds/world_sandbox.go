package worlds

import (
	"kosh/vpaul/floating/actors"
	"kosh/vpaul/floating/core"
)

func NewWorldSandBox() (*core.Actor, *core.World) {
	player := actors.NewCircle(400, 300)
	world := core.NewWorld()
	world.AddActor(player)
	// Static circles
	world.AddStaticActor(actors.NewStatic(500, 300))
	world.AddStaticActor(actors.NewStatic(250, 380))
	world.AddStaticActor(actors.NewStatic(620, 200))
	// Axis-aligned boxes
	world.AddStaticActor(actors.NewStaticSquare(350, 200, 55, 25))
	world.AddStaticActor(actors.NewStaticSquare(180, 260, 30, 50))
	// Rotated boxes
	world.AddStaticActor(actors.NewStaticRotatedSquare(480, 430, 60, 22, 30))
	world.AddStaticActor(actors.NewStaticRotatedSquare(300, 460, 50, 20, 45))
	world.AddStaticActor(actors.NewStaticRotatedSquare(650, 380, 45, 18, 60))
	world.AddStaticActor(actors.NewStaticRotatedSquare(150, 470, 40, 22, -36))
	world.AddStaticActor(actors.NewStaticRotatedSquare(560, 150, 35, 35, 22.5))
	world.AddStaticActor(actors.NewStaticRotatedSquare(700, 460, 55, 18, -45))
	world.AddStaticActor(actors.NewStaticRotatedSquare(100, 600, 555, 20, 0))

	return player, world
}
