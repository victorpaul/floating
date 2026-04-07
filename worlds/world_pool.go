package worlds

import (
	"kosh/vpaul/floating/actors"
	"kosh/vpaul/floating/core"
)

func NewWorldPool() (*core.Actor, *core.World) {
	player := actors.NewCircle(400, 300)
	//
	w := core.NewWorld()
	w.AddActor(player)

	var o float64 = 30
	w.AddStaticActor(actors.NewStaticSquare(o*5, o*10, o, o*2))
	w.AddStaticActor(actors.NewStaticSquare(o*9, o*12, o*5, o))

	//w.AddStaticActor(actors.NewStaticSquare(o*13, o*10, o, o*3))
	w.AddStaticActor(actors.NewStaticSquare(o*13, o*14, o, o*3))

	w.AddStaticActor(actors.NewStaticSquare(o*17, o*15, o*5, o))
	w.AddStaticActor(actors.NewStaticSquare(o*21, o*11, o, o*3))

	w.AddStaticActor(actors.NewStaticSquare(400, 450, 350, 5))

	return player, w
}
