package water

import (
	"math"
	"math/rand"

	"kosh/vpaul/floating/core"
	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// CellSize is the side length of one CA cell in world units.
// Smaller = more detail and smoother flow, higher CPU cost.
const CellSize = 4.0

type CellType uint8

const (
	CellEmpty CellType = iota
	CellSolid
	CellWater
)

// Grid is a fixed-size cellular automaton that simulates water.
// All geometry must be axis-aligned (no rotated boxes).
// Dynamic actors are not baked — only static geometry.
type Grid struct {
	cells []CellType
	moved []bool // prevents a cell from being processed twice per tick
	cols  int
	rows  int
	origX float64
	origY float64

	img    *ebiten.Image
	pixBuf []byte

	leftToRight bool // sweep direction, alternates every tick
}

// NewGrid creates a water grid covering the world region
// [origX, origX+worldW) × [origY, origY+worldH).
func NewGrid(worldW, worldH, origX, origY float64) *Grid {
	cols := int(math.Ceil(worldW / CellSize))
	rows := int(math.Ceil(worldH / CellSize))
	return &Grid{
		cells:  make([]CellType, cols*rows),
		moved:  make([]bool, cols*rows),
		cols:   cols,
		rows:   rows,
		origX:  origX,
		origY:  origY,
		img:    ebiten.NewImage(cols, rows),
		pixBuf: make([]byte, cols*rows*4),
	}
}

// BakeGeometry marks cells covered by static box colliders as CellSolid.
// Call once after world construction, and again after punching holes.
func (g *Grid) BakeGeometry(actors []*core.Actor) {
	for i := range g.cells {
		if g.cells[i] == CellSolid {
			g.cells[i] = CellEmpty
		}
	}
	for _, a := range actors {
		if a.Collider == nil || a.Collider.Shape != core.ShapeBox {
			continue
		}
		half := a.Collider.HalfSize
		c0, r0 := g.worldToCell(a.Pos.X-half.X, a.Pos.Y-half.Y)
		c1, r1 := g.worldToCell(a.Pos.X+half.X, a.Pos.Y+half.Y)
		for r := r0; r <= r1; r++ {
			for c := c0; c <= c1; c++ {
				if g.inBounds(c, r) {
					g.cells[g.idx(c, r)] = CellSolid
				}
			}
		}
	}
}

// AddWater fills a square region around (worldX, worldY) with water.
// Only empty cells are filled.
func (g *Grid) AddWater(worldX, worldY float64, radius int) {
	col, row := g.worldToCell(worldX, worldY)
	for dr := -radius; dr <= radius; dr++ {
		for dc := -radius; dc <= radius; dc++ {
			c, r := col+dc, row+dr
			if g.inBounds(c, r) && g.cells[g.idx(c, r)] == CellEmpty {
				g.cells[g.idx(c, r)] = CellWater
			}
		}
	}
}

// PunchHole clears solid cells in a square region around (worldX, worldY),
// allowing water to flow through.
func (g *Grid) PunchHole(worldX, worldY float64, radius int) {
	col, row := g.worldToCell(worldX, worldY)
	for dr := -radius; dr <= radius; dr++ {
		for dc := -radius; dc <= radius; dc++ {
			c, r := col+dc, row+dr
			if g.inBounds(c, r) && g.cells[g.idx(c, r)] == CellSolid {
				g.cells[g.idx(c, r)] = CellEmpty
			}
		}
	}
}

// Update runs one CA tick. Call every game frame.
func (g *Grid) Update() {
	for i := range g.moved {
		g.moved[i] = false
	}

	// Sweep bottom-to-top so falling water is not processed twice per tick.
	// Alternate horizontal direction each tick to prevent directional bias.
	for row := g.rows - 1; row >= 0; row-- {
		if g.leftToRight {
			for col := 0; col < g.cols; col++ {
				g.updateCell(row, col)
			}
		} else {
			for col := g.cols - 1; col >= 0; col-- {
				g.updateCell(row, col)
			}
		}
	}
	g.leftToRight = !g.leftToRight
}

func (g *Grid) updateCell(row, col int) {
	i := g.idx(col, row)
	if g.cells[i] != CellWater || g.moved[i] {
		return
	}

	below := row + 1

	// 1. Fall straight down.
	if below < g.rows {
		bi := g.idx(col, below)
		if g.cells[bi] == CellEmpty {
			g.cells[i] = CellEmpty
			g.cells[bi] = CellWater
			g.moved[bi] = true
			return
		}
	}

	// Pick a random lateral preference; used for both diagonal and sideways.
	d1, d2 := -1, 1
	if rand.Intn(2) == 1 {
		d1, d2 = 1, -1
	}

	// 2. Fall diagonally down.
	for _, d := range [2]int{d1, d2} {
		nc := col + d
		if below < g.rows && g.inBounds(nc, below) {
			di := g.idx(nc, below)
			if g.cells[di] == CellEmpty && !g.moved[di] {
				g.cells[i] = CellEmpty
				g.cells[di] = CellWater
				g.moved[di] = true
				return
			}
		}
	}

	// 3. Spread sideways.
	for _, d := range [2]int{d1, d2} {
		nc := col + d
		if g.inBounds(nc, row) {
			si := g.idx(nc, row)
			if g.cells[si] == CellEmpty && !g.moved[si] {
				g.cells[i] = CellEmpty
				g.cells[si] = CellWater
				g.moved[si] = true
				return
			}
		}
	}
}

// Draw renders water cells to screen. Solid cells are transparent (actors draw
// their own geometry). Call after world.Draw so water overlays terrain.
func (g *Grid) Draw(screen *ebiten.Image, camera utils.Vec) {
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			base := (r*g.cols + c) * 4
			if g.cells[g.idx(c, r)] == CellWater {
				g.pixBuf[base] = 30
				g.pixBuf[base+1] = 100
				g.pixBuf[base+2] = 220
				g.pixBuf[base+3] = 200
			} else {
				g.pixBuf[base] = 0
				g.pixBuf[base+1] = 0
				g.pixBuf[base+2] = 0
				g.pixBuf[base+3] = 0
			}
		}
	}
	g.img.WritePixels(g.pixBuf)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(CellSize, CellSize)
	op.GeoM.Translate(g.origX-camera.X, g.origY-camera.Y)
	screen.DrawImage(g.img, op)
}

func (g *Grid) idx(col, row int) int { return row*g.cols + col }
func (g *Grid) inBounds(col, row int) bool {
	return col >= 0 && col < g.cols && row >= 0 && row < g.rows
}
func (g *Grid) worldToCell(wx, wy float64) (int, int) {
	return int(math.Floor((wx - g.origX) / CellSize)),
		int(math.Floor((wy - g.origY) / CellSize))
}
