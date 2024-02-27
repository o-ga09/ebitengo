package orange

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	OrangeImage  *ebiten.Image
	Orange1Image *ebiten.Image
	Orange2Image *ebiten.Image
	Orange3Image *ebiten.Image
	Orange4Image *ebiten.Image
	Orange5Image *ebiten.Image

	//go:embed assets/orange.png
	Orange_png []byte

	//go:embed assets/orange_1.png
	Orange1_png []byte

	//go:embed assets/orange_2.png
	Orange2_png []byte

	//go:embed assets/orange_3.png
	Orange3_png []byte

	//go:embed assets/orange_4.png
	Orange4_png []byte

	//go:embed assets/orange_5.png
	Orange5_png []byte
)
