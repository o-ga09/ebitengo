package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ORANGE = iota
	HASSAKU
	DEKOPON
	HYUGANATU
	SETOKA
	LEMON
	MIKAN
)

const (
	screenWidth  = 480
	screenHeight = 600

	gravity  = 0.98
	friction = 0.98
	bounce   = 0.3
	spring   = 0.4
)

var (
	orangeImage  *ebiten.Image
	orange1Image *ebiten.Image
	orange2Image *ebiten.Image
	orange3Image *ebiten.Image
	orange4Image *ebiten.Image
	orange5Image *ebiten.Image

	//go:embed img/assets/orange.png
	orange_png []byte

	//go:embed img/assets/orange_1.png
	orange1_png []byte

	//go:embed img/assets/orange_2.png
	orange2_png []byte

	//go:embed img/assets/orange_3.png
	orange3_png []byte

	//go:embed img/assets/orange_4.png
	orange4_png []byte

	//go:embed img/assets/orange_5.png
	orange5_png []byte
)

type Orange struct {
	X      float64
	Y      float64
	VX     float64
	VY     float64
	Radius float64
	Type   int
	Remove bool
}

func NewOrange(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 20,
		Type:   ORANGE,
	}
}

func NewHassaku(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 35,
		Type:   HASSAKU,
	}
}

func NewDekopon(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 50,
		Type:   DEKOPON,
	}
}

func NewHyuganatu(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 65,
		Type:   HYUGANATU,
	}
}

func NewLemon(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 80,
		Type:   LEMON,
	}
}

func NewMikan(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 20,
		Type:   MIKAN,
	}
}

type Draw struct {
	op ebiten.DrawImageOptions
}

func init() {
	orangeImage = loadImage(orange_png)
	orange1Image = loadImage(orange1_png)
	orange2Image = loadImage(orange2_png)
	orange3Image = loadImage(orange3_png)
	orange4Image = loadImage(orange4_png)
	orange5Image = loadImage(orange5_png)
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
	var img *ebiten.Image
	switch {
	case orange.Type == ORANGE:
		img = orangeImage
	case orange.Type == HASSAKU:
		img = orange1Image
	case orange.Type == DEKOPON:
		img = orange2Image
	case orange.Type == HYUGANATU:
		img = orange3Image
	case orange.Type == LEMON:
		img = orange4Image
	case orange.Type == MIKAN:
		img = orange5Image
	}

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
		NewOrange(100, 0),
		NewOrange(110, -100),
		// NewDekopon(100, -1000),
		// NewHassaku(110, -2000),
		// NewHyuganatu(120, -3000),
		// NewLemon(150, -4000),
		// NewMikan(170, -5000),
	}

	world = World{X: 0, Y: 0, Width: screenWidth, Height: screenHeight}
	draw  = &Draw{}
	calc  = &Calc{World: world}
)

type Calc struct {
	World World
	Score int
}

func (c *Calc) Oranges(oranges []*Orange) []*Orange {
	oranges = c.combine(oranges)
	c.move(oranges)
	c.Collision(oranges)
	c.screenWrap(oranges)

	return oranges
}

func (c *Calc) move(o []*Orange) {
	for _, orange := range oranges {
		orange.VX *= friction
		orange.VY *= friction
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

func (c *Calc) Collision(oranges []*Orange) {
	for i, orange := range oranges {
		for j := i + 1; j < len(oranges); j++ {
			g := oranges[j]
			dx := g.X - orange.X
			dy := g.Y - orange.Y
			d := math.Sqrt(dx*dx + dy*dy)
			minD := orange.Radius + g.Radius
			if d < minD {
				angle := math.Atan2(dy, dx)
				tx := orange.X + math.Cos(angle)*minD
				ty := orange.Y + math.Sin(angle)*minD
				ax := (tx - g.X) * spring
				ay := (ty - g.Y) * spring
				orange.VX -= ax
				orange.VY -= ay
				g.VX += ax
				g.VY += ay

				orange.X = orange.X - math.Cos(angle)*(minD-d)/2
				orange.Y = orange.Y - math.Sin(angle)*(minD-d)/2
				g.X = g.X + math.Cos(angle)*(minD-d)/2
				g.Y = g.Y + math.Sin(angle)*(minD-d)/2
			}
		}
	}
}

func (c *Calc) combine(oranges []*Orange) []*Orange {
	newOranges := make([]*Orange, 0)
	for i, orange := range oranges {
		for j := i + 1; j < len(oranges); j++ {
			g := oranges[j]
			if orange.Remove || g.Remove {
				continue
			}

			dx := g.X - orange.X
			dy := g.Y - orange.Y
			d := math.Sqrt(dx*dx + dy*dy)
			minD := orange.Radius + g.Radius
			if int64(d) <= int64(minD) && orange.Type == g.Type {
				orange.Remove = true
				g.Remove = true
				var next *Orange
				switch orange.Type {
				case ORANGE:
					next = NewHassaku((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 10
				case HASSAKU:
					next = NewDekopon((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 20
				case DEKOPON:
					next = NewHyuganatu((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 30
				case HYUGANATU:
					next = NewLemon((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 40
				case LEMON:
					next = NewMikan((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 50
				case MIKAN:
					next = NewMikan((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 60
				}
				if next != nil {
					newOranges = append(newOranges, next)
				}
			}
		}
	}

	for _, orange := range oranges {
		if !orange.Remove {
			newOranges = append(newOranges, orange)
		}
	}

	return newOranges
}

type World struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type Game struct{}

func (g *Game) Update() error {
	oranges = calc.Oranges(oranges)
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
