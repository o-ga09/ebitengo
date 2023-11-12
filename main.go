package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	runnerImage *ebiten.Image
	posX        int = 0
	posY        int = 0
	runes       []rune
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	count int
}

func NewGame() *Game {
	return &Game{count: 0}
}

func (g *Game) Update() error {
	// 押されていたフレームの間trueになる
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		// キーボードのAが押された場合の処理
		fmt.Println("A is Pushed !")
	}
	// 押されたフレームのみtrueになる
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		// キーボードのAが押された場合の処理
		fmt.Println("just A is Pushed !")
	}

	// キーを離したフレームのみtrueになる
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		// キーボードのAが離された場合の処理
		fmt.Println("A is released !")
	}

	// 左クリックした場合trueになる
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// マウスの左ボタンが押された場合の処理
		fmt.Println("LeftClicked !")
	}

	// 右クリックした場合trueになる
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		// マウスの右ボタンが押された場合の処理
		fmt.Println("RightClicked !")
	}

	// 押されたフレームのみtrueになる
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// キーボードのAが押された場合の処理
		fmt.Println("just Left is Clicled !")
	}

	// キーを離したフレームのみtrueになる
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		// キーボードのAが離された場合の処理
		fmt.Println("just Right is Clicked !")
	}

	posX, posY = ebiten.CursorPosition()
	runes = ebiten.AppendInputChars(runes)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(float64(posX), float64(posY))
	screen.DrawImage(runnerImage, op)

	// runesに入力された文字情報があるので下記のように出力することが可能
	ebitenutil.DebugPrint(screen, string(runes))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X : %d, Y : %d", posX, posY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// ウィンドウサイズの指定
	ebiten.SetWindowSize(640, 480)
	// ウィンドウタイトルの指定
	ebiten.SetWindowTitle("Hello, World!")

	// 画像ファイルを開く
	file, err := os.Open("img/gopher.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 画像を読み込む
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	runnerImage = ebiten.NewImageFromImage(img)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
