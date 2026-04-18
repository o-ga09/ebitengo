package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	game "github.com/o-ga09/ebitengo/internal"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Wave Duel")

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
