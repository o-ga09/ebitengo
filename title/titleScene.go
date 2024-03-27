package title

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

var (
	TitleBackgroundImage *ebiten.Image
	StartImage           *ebiten.Image
	AlertImage           *ebiten.Image

	//go:embed assets/title.png
	TitleImage_png []byte

	//go:embed assets/start.png
	StartImage_png []byte

	//go:embed assets/alert.png
	AlertImage_png []byte
)

type TitleScene struct {
	op ebiten.DrawImageOptions
}

func (d *TitleScene) TitleScreen(screen *ebiten.Image) {
	Startimg := StartImage

	w, h := Startimg.Bounds().Dx(), Startimg.Bounds().Dy()

	// 描画位置をリセット
	d.op.GeoM.Reset()

	// 画面サイズに合わせるようにリサイズ
	d.op.GeoM.Scale(float64(ScreenWidth)/float64(w)*0.5, float64(ScreenHeight)/float64(h)*0.25)
	d.op.GeoM.Translate(float64(ScreenWidth/4), float64(ScreenHeight/2))

	// 画面に描画
	screen.DrawImage(Startimg, &d.op)
}

func (d *TitleScene) AlertScreen(screen *ebiten.Image) {
	Alertimg := AlertImage

	w, h := Alertimg.Bounds().Dx(), Alertimg.Bounds().Dy()

	// 描画位置をリセット
	d.op.GeoM.Reset()

	// 画面サイズに合わせるようにリサイズ
	d.op.GeoM.Scale(float64(ScreenWidth)/float64(w)*0.5, float64(ScreenHeight)/float64(h)*0.25)
	d.op.GeoM.Translate(float64(ScreenWidth/4), float64(ScreenHeight/8))

	// 画面に描画
	screen.DrawImage(Alertimg, &d.op)
}
