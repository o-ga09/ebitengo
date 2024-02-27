package main

import (
	_ "embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	game "github.com/o-ga09/ebiten/Game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("オレンジ")
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		panic(err)
	}
}
