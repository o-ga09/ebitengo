package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600

	ampStep   = 10.0        // 振幅変化量 [px/frame]
	freqStep  = 0.05        // 周波数変化量 [Hz/frame]
	phaseStep = math.Pi / 8 // 位相変化量 [rad/frame]
	ampMin    = 10.0
	ampMax    = 240.0
	freqMin   = 0.1
	freqMax   = 5.0
)

// playerKeys は1人分のキーバインドを保持する。
type playerKeys struct {
	ampUp, ampDown     ebiten.Key
	freqUp, freqDown   ebiten.Key
	phaseUp, phaseDown ebiten.Key
	fire               ebiten.Key
	waveTypes          [6]ebiten.Key
}

var (
	p1Keys = playerKeys{
		ampUp: ebiten.KeyW, ampDown: ebiten.KeyS,
		freqUp: ebiten.KeyE, freqDown: ebiten.KeyQ,
		phaseUp: ebiten.KeyX, phaseDown: ebiten.KeyZ,
		fire:      ebiten.KeySpace,
		waveTypes: [6]ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5, ebiten.Key6},
	}
	p2Keys = playerKeys{
		ampUp: ebiten.KeyArrowUp, ampDown: ebiten.KeyArrowDown,
		freqUp: ebiten.KeyPeriod, freqDown: ebiten.KeyComma,
		phaseUp: ebiten.KeyM, phaseDown: ebiten.KeyN,
		fire:      ebiten.KeyEnter,
		waveTypes: [6]ebiten.Key{ebiten.KeyNumpad1, ebiten.KeyNumpad2, ebiten.KeyNumpad3, ebiten.KeyNumpad4, ebiten.KeyNumpad5, ebiten.KeyNumpad6},
	}
)

// Player は1人のプレイヤー状態を保持する。
type Player struct {
	pending     Wave
	fired       []Wave
	firePressed bool
	waveColor   color.NRGBA // プレーヤー識別色
}

// Game はゲーム全体の状態を管理する。
type Game struct {
	tick    int
	players [2]*Player
}

// NewGame はゲームを初期化して返す。
func NewGame() *Game {
	return &Game{
		players: [2]*Player{
			{
				pending:   Wave{Amplitude: 80, Frequency: 0.5, Phase: 0, Type: WaveSin, Direction: 1},
				waveColor: color.NRGBA{R: 0, G: 220, B: 255, A: 200},
			},
			{
				pending:   Wave{Amplitude: 80, Frequency: 0.5, Phase: 0, Type: WaveSin, Direction: -1},
				waveColor: color.NRGBA{R: 255, G: 140, B: 0, A: 200},
			},
		},
	}
}

func updatePlayer(p *Player, keys playerKeys) {
	waveTypeList := []WaveType{WaveSin, WaveCos, WaveTan, WaveSquare, WaveSawtooth, WaveTriangle}
	for i, k := range keys.waveTypes {
		if ebiten.IsKeyPressed(k) {
			p.pending.Type = waveTypeList[i]
		}
	}
	if ebiten.IsKeyPressed(keys.ampUp) {
		p.pending.Amplitude = math.Min(p.pending.Amplitude+ampStep, ampMax)
	}
	if ebiten.IsKeyPressed(keys.ampDown) {
		p.pending.Amplitude = math.Max(p.pending.Amplitude-ampStep, ampMin)
	}
	if ebiten.IsKeyPressed(keys.freqUp) {
		p.pending.Frequency = math.Min(p.pending.Frequency+freqStep, freqMax)
	}
	if ebiten.IsKeyPressed(keys.freqDown) {
		p.pending.Frequency = math.Max(p.pending.Frequency-freqStep, freqMin)
	}
	if ebiten.IsKeyPressed(keys.phaseUp) {
		p.pending.Phase = math.Mod(p.pending.Phase+phaseStep, 2*math.Pi)
	}
	if ebiten.IsKeyPressed(keys.phaseDown) {
		p.pending.Phase = math.Mod(p.pending.Phase-phaseStep+2*math.Pi, 2*math.Pi)
	}

	fireNow := ebiten.IsKeyPressed(keys.fire)
	if fireNow && !p.firePressed {
		fired := p.pending
		p.fired = append(p.fired, fired)
	}
	p.firePressed = fireNow
}

// Update はフレームごとのロジック更新（毎秒60回呼ばれる）。
func (g *Game) Update() error {
	g.tick++
	updatePlayer(g.players[0], p1Keys)
	updatePlayer(g.players[1], p2Keys)

	// R でリセット
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.players[0].fired = nil
		g.players[1].fired = nil
	}

	return nil
}

// Draw はフレームごとの描画処理。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 15, G: 15, B: 30, A: 255})

	t := float64(g.tick) / 60.0
	baseline := float64(screenHeight) / 2

	// xMin〜xMax の範囲だけ折れ線で描画
	drawWaveRange := func(waves []Wave, col color.NRGBA, thick float32, xMin, xMax int) {
		for x := xMin + 1; x <= xMax; x++ {
			y0 := baseline - CompositeAmplitude(waves, float64(x-1), t)
			y1 := baseline - CompositeAmplitude(waves, float64(x), t)
			vector.StrokeLine(screen, float32(x-1), float32(y0), float32(x), float32(y1), thick, col, true)
		}
	}

	// P1 pending プレビュー（自分の左半分のみ、薄い・静止）
	p1Preview := g.players[0].pending
	p1Preview.Frequency = 0 // 静止プレビュー
	drawWaveRange([]Wave{p1Preview},
		color.NRGBA{R: g.players[0].waveColor.R, G: g.players[0].waveColor.G, B: g.players[0].waveColor.B, A: 60},
		1, 0, screenWidth/2)

	// P2 pending プレビュー（自分の右半分のみ、薄い・静止）
	p2Preview := g.players[1].pending
	p2Preview.Frequency = 0
	drawWaveRange([]Wave{p2Preview},
		color.NRGBA{R: g.players[1].waveColor.R, G: g.players[1].waveColor.G, B: g.players[1].waveColor.B, A: 60},
		1, screenWidth/2, screenWidth-1)

	// 発射済み波の合成波（白、太め）―発射がある場合のみ
	var allFired []Wave
	for _, p := range g.players {
		allFired = append(allFired, p.fired...)
	}
	if len(allFired) > 0 {
		drawWaveRange(allFired, color.NRGBA{R: 255, G: 255, B: 220, A: 255}, 3, 0, screenWidth-1)
	}

	// 基準線（水平）
	vector.StrokeLine(screen, 0, float32(baseline), screenWidth, float32(baseline), 1,
		color.NRGBA{R: 60, G: 60, B: 80, A: 180}, false)

	// 中央の区切り線
	vector.StrokeLine(screen, screenWidth/2, 0, screenWidth/2, screenHeight, 1,
		color.NRGBA{R: 100, G: 100, B: 100, A: 100}, false)

	// --- HUD ---
	p1 := g.players[0].pending
	p2 := g.players[1].pending
	ebitenutil.DebugPrintAt(screen, "--- P1 (left→right) ---", 8, 8)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type:%-8s  Amp:%.0f  Freq:%.2f  Phase:%.2fπ",
		p1.Type, p1.Amplitude, p1.Frequency, p1.Phase/math.Pi), 8, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[1-6]Type [W/S]Amp [Q/E]Freq [Z/X]Phase [Space]Fire  Fired:%d", len(g.players[0].fired)), 8, 32)

	p2info := fmt.Sprintf("Type:%-8s  Amp:%.0f  Freq:%.2f  Phase:%.2fπ", p2.Type, p2.Amplitude, p2.Frequency, p2.Phase/math.Pi)
	ebitenutil.DebugPrintAt(screen, "--- P2 (right←left) ---", screenWidth/2+8, 8)
	ebitenutil.DebugPrintAt(screen, p2info, screenWidth/2+8, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[Num1-6]Type [↑↓]Amp [,.]Freq [N/M]Phase [Enter]Fire  Fired:%d", len(g.players[1].fired)), screenWidth/2+8, 32)
	// リセット
	ebitenutil.DebugPrintAt(screen, "[R] Reset all", screenWidth/2-40, screenHeight-16)
}

// Layout は論理画面サイズを返す。
func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}
