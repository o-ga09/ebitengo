package game

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	orange "github.com/o-ga09/ebiten/Orange"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

var (
	backgroundImage *ebiten.Image

	//go:embed assets/background.png
	background_png []byte

	world = orange.World{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight}
	draw  = &orange.Draw{}
	calc  = &orange.Calc{World: world}
)

type Game struct{}

func init() {
	orange.OrangeImage = orange.LoadImage(orange.Orange_png)
	orange.Orange1Image = orange.LoadImage(orange.Orange1_png)
	orange.Orange2Image = orange.LoadImage(orange.Orange2_png)
	orange.Orange3Image = orange.LoadImage(orange.Orange3_png)
	orange.Orange4Image = orange.LoadImage(orange.Orange4_png)
	orange.Orange5Image = orange.LoadImage(orange.Orange5_png)
}

func (g *Game) Update() error {
	orange.Oranges = calc.Oranges(orange.Oranges)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw.World(screen, world)
	draw.Oranges(screen, world, orange.Oranges)
}

func (g *Game) Layout(outsideWidth, outsodeHight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
