package worlds

import (
	"kosh/vpaul/floating/actors"
	"kosh/vpaul/floating/core"
)

func NewWorld() (*core.Actor, *core.World) {
	player := actors.NewCircle(400, 300)

	w := core.NewWorld()

	return player, w
}
