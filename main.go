package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 600

	gravity = 0.4
	bounce  = 1.0
)

var (
	orangeImage *ebiten.Image

	//go:embed img/assets/orange.png
	orange_png []byte
)

type Orange struct {
	X      float64
	Y      float64
	VX     float64
	VY     float64
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

func (d *Draw) Oranges(screen *ebiten.Image, world World, oranges []*Orange) {
	for _, o := range oranges {
		d.Orange(screen, world, o)
	}
}
func (d *Draw) World(screen *ebiten.Image, world World) {
	vector.DrawFilledRect(screen, float32(world.X), float32(world.Y), float32(world.Width), float32(world.Height), color.RGBA{0x66, 0x66, 0x66, 0xff}, false)
}

func (d *Draw) Orange(screen *ebiten.Image, world World, orange *Orange) {
	img := orangeImage

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	d.op.Filter = ebiten.FilterLinear
	d.op.GeoM.Reset()
	d.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	d.op.GeoM.Scale(orange.Radius/float64(w)*2, orange.Radius/float64(h)*2)
	d.op.GeoM.Translate(float64(world.X), float64(world.Y))
	d.op.GeoM.Translate(float64(orange.X), float64(orange.Y))
	screen.DrawImage(img, &d.op)
}

var (
	oranges = []*Orange{
		{X: 100, Y: 100, VX: 15, VY: 2, Radius: 25},
		{X: 250, Y: 100, VX: -15, VY: 4, Radius: 50},
	}

	world = World{X: 0, Y: 0, Width: screenWidth, Height: screenHeight}
	draw  = &Draw{}
	calc  = &Calc{World: world}
)

type Calc struct {
	World World
}

func (c *Calc) Oranges(oranges []*Orange) {
	c.move(oranges)
	c.screenWrap(oranges)
}

func (c *Calc) move(o []*Orange) {
	for _, orange := range oranges {
		orange.VY += gravity
		orange.X += orange.VX
		orange.Y += orange.VY
	}
}

func (c *Calc) screenWrap(oranges []*Orange) {
	for _, orange := range oranges {
		// X軸方向にはみ出す
		if orange.X-orange.Radius < 0 {
			orange.X = orange.Radius
			orange.VX *= -bounce
		} else if c.World.Width < orange.X+orange.Radius {
			orange.X = c.World.Width - orange.Radius
			orange.VX *= -bounce
		}

		// Y軸方向にはみ出す
		if orange.Y < 0 {
			// なにもしない
		} else if c.World.Height < orange.Y+orange.Radius {
			orange.Y = c.World.Height - orange.Radius
			orange.VY *= -bounce
		}
	}
}

type World struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type Game struct{}

func (g *Game) Update() error {
	calc.Oranges(oranges)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw.World(screen, world)
	draw.Oranges(screen, world, oranges)
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
