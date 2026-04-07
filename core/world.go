package core

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// World owns all actors, including static geometry.
type World struct {
	Actors    []*Actor
	Grid      *SpatialGrid
	DebugGrid bool
}

func NewWorld() *World {
	return &World{Grid: NewSpatialGrid()}
}

// AddActor registers a dynamic actor (tracked in the grid, updated every frame).
func (w *World) AddActor(a *Actor) {
	w.Actors = append(w.Actors, a)
	w.Grid.Register(a)
}

// AddStaticActor registers a static actor (placed in all cells its AABB covers,
// never moved in the grid).
func (w *World) AddStaticActor(a *Actor) {
	w.Actors = append(w.Actors, a)
	w.Grid.RegisterStatic(a)
}

// RemoveActor removes a dynamic actor from the world and the grid.
func (w *World) RemoveActor(a *Actor) {
	w.Grid.Unregister(a)
	for i, actor := range w.Actors {
		if actor == a {
			w.Actors[i] = w.Actors[len(w.Actors)-1]
			w.Actors = w.Actors[:len(w.Actors)-1]
			return
		}
	}
}

func (w *World) Update(ctx *GameContext) {
	for _, a := range w.Actors {
		a.Update(ctx)
		w.Grid.UpdateActor(a)
	}
}

func (w *World) Draw(screen *ebiten.Image, ctx *GameContext) {
	if w.DebugGrid {
		w.drawDebugGrid(screen, ctx)
	}
	for _, a := range w.Actors {
		a.Render(screen, ctx)
	}
}

func (w *World) drawDebugGrid(screen *ebiten.Image, ctx *GameContext) {
	const size = float32(SpatialGridCellSize)
	cellColor := color.RGBA{R: 80, G: 180, B: 80, A: 120}
	textColor := color.RGBA{R: 180, G: 255, B: 180, A: 255}
	_ = textColor // ebitenutil.DebugPrintAt uses its own colour; kept for reference

	for key, actors := range w.Grid.cells {
		if len(actors) == 0 {
			continue
		}
		col := int64(int32(key & 0xFFFFFFFF))
		row := key >> 32

		wx := float32(col)*size - float32(ctx.Camera.X)
		wy := float32(row)*size - float32(ctx.Camera.Y)

		vector.StrokeRect(screen, wx, wy, size, size, 1, cellColor, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", len(actors)),
			int(wx)+4, int(wy)+4)
	}
}
