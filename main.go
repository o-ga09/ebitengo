package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 480
	screenHeight = 600
)

var (
	orangeImage *ebiten.Image

	//go:embed img/assets/orange.png
	orange_png []byte
)

type Orange struct {
	X      float64
	Y      float64
	Radius float64
}

type Draw struct {
	op ebiten.DrawImageOptions
}

func init() {
	orangeImage = loadImage(orange_png)
}

func loadImage(b []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	origin := ebiten.NewImageFromImage(img)
	s := origin.Bounds().Size()
	ebitenImg := ebiten.NewImage(s.X, s.Y)
	op := &ebiten.DrawImageOptions{}
	ebitenImg.DrawImage(origin, op)
	return ebitenImg
}

func (d *Draw) Oranges(screen *ebiten.Image, oranges []*Orange) {
	for _, o := range oranges {
		d.Orange(screen, o)
	}
}

func (d *Draw) Orange(screen *ebiten.Image, orange *Orange) {
	img := orangeImage

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	d.op.Filter = ebiten.FilterLinear
	d.op.GeoM.Reset()
	d.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	d.op.GeoM.Scale(orange.Radius/float64(w)*2, orange.Radius/float64(h)*2)
	d.op.GeoM.Translate(float64(orange.X), float64(orange.Y))
	screen.DrawImage(img, &d.op)
}

var (
	oranges = []*Orange{
		{X: 100, Y: 100, Radius: 25},
		{X: 250, Y: 200, Radius: 50},
	}
	draw = &Draw{}
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw.Oranges(screen, oranges)
}

func (g *Game) Layout(outsideWidth, outsodeHight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("オレンジ")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
