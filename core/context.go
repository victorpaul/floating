package core

import "kosh/vpaul/floating/utils"

// GameContext is passed to every Update and Draw call,
// giving any subsystem read access to the full game state for the current frame.
type GameContext struct {
	Camera     utils.Vec
	World      *World
	Controller *Controller
}
