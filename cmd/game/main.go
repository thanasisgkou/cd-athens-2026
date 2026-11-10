package main

import (
	"log"

	"athens-rampage/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Athens Rampage: Τιτάνες στην Ακρόπολη")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g, err := game.NewAthensGame()
	if err != nil {
		log.Fatalf("Failed to initialize Athens Rampage: %v", err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game loop terminated: %v", err)
	}
}
