package core

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DebugOverlay struct {
	updateDuration time.Duration
	drawDuration   time.Duration
	ActorCount     int
}

func NewDebugOverlay() *DebugOverlay {
	return &DebugOverlay{}
}

func (d *DebugOverlay) BeginUpdate() time.Time {
	return time.Now()
}

func (d *DebugOverlay) EndUpdate(start time.Time) {
	d.updateDuration = time.Since(start)
}

func (d *DebugOverlay) BeginDraw() time.Time {
	return time.Now()
}

func (d *DebugOverlay) EndDraw(start time.Time) {
	d.drawDuration = time.Since(start)
}

func (d *DebugOverlay) Draw(screen *ebiten.Image) {
	updateMs := float64(d.updateDuration.Microseconds()) / 1000.0
	drawMs := float64(d.drawDuration.Microseconds()) / 1000.0
	totalMs := updateMs + drawMs

	maxFPS := 0.0
	if totalMs > 0 {
		maxFPS = 1000.0 / totalMs
	}

	msg := fmt.Sprintf(
		"FPS: %.0f\nMax FPS: %.0f\nUpdate: %.2f ms\nDraw: %.2f ms\nActors: %d",
		ebiten.ActualFPS(),
		maxFPS,
		updateMs,
		drawMs,
		d.ActorCount,
	)
	ebitenutil.DebugPrintAt(screen, msg, 4, 4)
}
