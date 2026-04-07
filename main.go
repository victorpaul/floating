package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(int(ScreenWidth), int(ScreenHeight))
	ebiten.SetWindowTitle("Floating")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGameScreen()); err != nil {
		log.Fatal(err)
	}
}
