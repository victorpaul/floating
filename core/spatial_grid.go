package core

import (
	"math"

	"kosh/vpaul/floating/utils"
)

// SpatialGridCellSize is the side length of each grid cell in world units.
// Tune this to roughly match the largest collider radius/half-extent in the scene.
const SpatialGridCellSize = 100.0

// SpatialGrid partitions actors into fixed-size cells so collision checks only
// need to consider a small neighbourhood instead of every actor in the world.
//
//   - Dynamic actors (tracked in actorCells) are registered with Register and
//     must be refreshed every frame via UpdateActor.
//   - Static actors are registered with RegisterStatic across all cells their
//     AABB overlaps; they are never moved, so no entry is kept in actorCells.
type SpatialGrid struct {
	cells      map[int64][]*Actor
	actorCells map[*Actor]int64
}

func NewSpatialGrid() *SpatialGrid {
	return &SpatialGrid{
		cells:      make(map[int64][]*Actor),
		actorCells: make(map[*Actor]int64),
	}
}

// cellKey returns a unique int64 key for the cell that contains world position pos.
func cellKey(pos utils.Vec) int64 {
	col := int64(math.Floor(pos.X / SpatialGridCellSize))
	row := int64(math.Floor(pos.Y / SpatialGridCellSize))
	// Pack into int64 using the upper 32 bits for row, lower 32 for col.
	return row<<32 | (col & 0xFFFFFFFF)
}

// Register adds a dynamic actor to the grid. Call UpdateActor every frame to
// keep its cell current as it moves.
func (g *SpatialGrid) Register(a *Actor) {
	key := cellKey(a.Pos)
	g.cells[key] = append(g.cells[key], a)
	g.actorCells[a] = key
}

// RegisterStatic inserts a static actor into every cell its collider's AABB
// overlaps. The actor is never tracked for movement.
func (g *SpatialGrid) RegisterStatic(a *Actor) {
	var halfX, halfY float64
	if a.Collider != nil {
		switch a.Collider.Shape {
		case ShapeBox:
			cos := math.Abs(math.Cos(a.Collider.Angle))
			sin := math.Abs(math.Sin(a.Collider.Angle))
			halfX = a.Collider.HalfSize.X*cos + a.Collider.HalfSize.Y*sin
			halfY = a.Collider.HalfSize.X*sin + a.Collider.HalfSize.Y*cos
		case ShapeCircle:
			halfX = a.Collider.Radius
			halfY = a.Collider.Radius
		}
	}

	startCol := int64(math.Floor((a.Pos.X - halfX) / SpatialGridCellSize))
	endCol := int64(math.Floor((a.Pos.X + halfX) / SpatialGridCellSize))
	startRow := int64(math.Floor((a.Pos.Y - halfY) / SpatialGridCellSize))
	endRow := int64(math.Floor((a.Pos.Y + halfY) / SpatialGridCellSize))

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			key := row<<32 | (col & 0xFFFFFFFF)
			g.cells[key] = append(g.cells[key], a)
		}
	}
}

// Unregister removes a dynamic actor from the grid.
func (g *SpatialGrid) Unregister(a *Actor) {
	key, ok := g.actorCells[a]
	if !ok {
		return
	}
	g.removeFromCell(key, a)
	delete(g.actorCells, a)
}

// UpdateActor moves a dynamic actor to its new cell if its position has crossed
// a cell boundary. Safe to call every frame; it is a no-op when the cell is unchanged
// or when the actor is not tracked (e.g. statics).
func (g *SpatialGrid) UpdateActor(a *Actor) {
	oldKey, ok := g.actorCells[a]
	if !ok {
		return
	}
	newKey := cellKey(a.Pos)
	if oldKey == newKey {
		return
	}
	g.removeFromCell(oldKey, a)
	g.cells[newKey] = append(g.cells[newKey], a)
	g.actorCells[a] = newKey
}

// GetNearby returns all actors in the 3×3 cell neighbourhood around pos.
// The returned slice may contain duplicates if a static actor spans multiple cells.
func (g *SpatialGrid) GetNearby(pos utils.Vec) []*Actor {
	col := int64(math.Floor(pos.X / SpatialGridCellSize))
	row := int64(math.Floor(pos.Y / SpatialGridCellSize))

	var result []*Actor
	for dr := int64(-1); dr <= 1; dr++ {
		for dc := int64(-1); dc <= 1; dc++ {
			key := (row+dr)<<32 | ((col + dc) & 0xFFFFFFFF)
			result = append(result, g.cells[key]...)
		}
	}
	return result
}

func (g *SpatialGrid) removeFromCell(key int64, a *Actor) {
	list := g.cells[key]
	for i, u := range list {
		if u == a {
			list[i] = list[len(list)-1]
			g.cells[key] = list[:len(list)-1]
			return
		}
	}
}
