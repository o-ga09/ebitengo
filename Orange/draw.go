package orange

import (
	"bytes"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Draw struct {
	op ebiten.DrawImageOptions
}

func LoadImage(b []byte) *ebiten.Image {
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
		img = OrangeImage
	case orange.Type == HASSAKU:
		img = Orange1Image
	case orange.Type == DEKOPON:
		img = Orange2Image
	case orange.Type == HYUGANATU:
		img = Orange3Image
	case orange.Type == LEMON:
		img = Orange4Image
	case orange.Type == MIKAN:
		img = Orange5Image
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
