package pkg

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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
