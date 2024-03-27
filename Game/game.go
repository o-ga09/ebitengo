package game

import (
	_ "embed"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	orange "github.com/o-ga09/ebiten/Orange"
	"github.com/o-ga09/ebiten/pkg"
	"github.com/o-ga09/ebiten/title"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth      = 640
	ScreenHeight     = 480
	TITILE       int = iota
	GAME
	FIN
)

var (
	world = orange.World{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight}
	draw  = &orange.Draw{}
	calc  = &orange.Calc{World: world}

	// title Scene
	titleScene = &title.TitleScene{}

	mplusBigFont font.Face
)

type Game struct {
	DummyFlg bool
	Scene    int
}

func init() {
	title.TitleBackgroundImage = pkg.LoadImage(title.TitleImage_png)
	title.StartImage = pkg.LoadImage(title.StartImage_png)
	title.AlertImage = pkg.LoadImage(title.AlertImage_png)

	orange.OrangeImage = pkg.LoadImage(orange.Orange_png)
	orange.Orange1Image = pkg.LoadImage(orange.Orange1_png)
	orange.Orange2Image = pkg.LoadImage(orange.Orange2_png)
	orange.Orange3Image = pkg.LoadImage(orange.Orange3_png)
	orange.Orange4Image = pkg.LoadImage(orange.Orange4_png)
	orange.Orange5Image = pkg.LoadImage(orange.Orange5_png)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont = text.FaceWithLineHeight(mplusBigFont, 54)
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	// シーンの切り替え
	switch g.Scene {
	case TITILE:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Scene = GAME
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && ((x > 189 && x < 454) && (y > 243 && y < 360)) {
			g.DummyFlg = true
		}
	case GAME:
		if ebiten.IsKeyPressed(ebiten.KeyO) {
			g.Scene = FIN
		}
		orange.Oranges = calc.Oranges(orange.Oranges)
	case FIN:
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.Scene = TITILE
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// シーン毎に描画する対象を切り替える
	switch g.Scene {
	case TITILE:
		ebitenutil.DebugPrint(screen, "title")
		titleScene.TitleScreen(screen)

		if g.DummyFlg {
			titleScene.AlertScreen(screen)
		}
	case GAME:
		ebitenutil.DebugPrint(screen, "in game")
		draw.World(screen, world)
		draw.Oranges(screen, world, orange.Oranges)
	case FIN:
		ebitenutil.DebugPrint(screen, "finish")
		text.Draw(screen, "終了！！！\n\n Escape: ゲームを終了 \n E: タイトルに戻る", mplusBigFont, ScreenWidth/8, ScreenHeight/2, color.RGBA{R: 255, G: 255, B: 0, A: 0})
	}
}

func (g *Game) Layout(outsideWidth, outsodeHight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
