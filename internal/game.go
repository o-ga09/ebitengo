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

	resolveTicks = 120 // 解決フェーズ: 2秒
	damageTicks  = 90  // ダメージ表示: 1.5秒

	waveBaseline = 310 // 波形の基準 y 座標
)

// BattleState はバトルの進行状態を表す。
type BattleState int

const (
	StateTurnP1     BattleState = iota // P1 が調整・発射するターン
	StateTurnP2                        // P2 が調整・発射するターン
	StateResolve                       // 合成波アニメーション（最大振幅を計測）
	StateShowDamage                    // ダメージ表示
	StateGameOver                      // 勝敗決定
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
	waveColor   color.NRGBA
}

// Game はゲーム全体の状態を管理する。
type Game struct {
	tick            int
	players         [2]*Player
	hp              [2]int
	battleState     BattleState
	phaseTick       int
	maxAmpPerPlayer [2]float64 // 各プレイヤーの波だけの最大振幅
	lastDamage      [2]int
	winner          int // 0=P1勝利, 1=P2勝利, -1=引き分け
}

// NewGame はゲームを初期化して返す。
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.tick = 0
	g.hp = [2]int{InitialHP, InitialHP}
	g.battleState = StateTurnP1
	g.phaseTick = 0
	g.maxAmpPerPlayer = [2]float64{}
	g.lastDamage = [2]int{}
	g.winner = 0
	g.players = [2]*Player{
		{
			pending:   Wave{Amplitude: 80, Frequency: 0.5, Phase: 0, Type: WaveSin, Direction: 1},
			waveColor: color.NRGBA{R: 0, G: 220, B: 255, A: 200},
		},
		{
			pending:   Wave{Amplitude: 80, Frequency: 0.5, Phase: 0, Type: WaveSin, Direction: -1},
			waveColor: color.NRGBA{R: 255, G: 140, B: 0, A: 200},
		},
	}
}

// updatePlayerTurn は現在のターンプレイヤーの操作を処理し、発射した場合 true を返す。
func updatePlayerTurn(p *Player, keys playerKeys) bool {
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
	fired := fireNow && !p.firePressed
	if fired {
		p.fired = append(p.fired, p.pending)
	}
	p.firePressed = fireNow
	return fired
}

// Update はフレームごとのロジック更新（毎秒60回呼ばれる）。
func (g *Game) Update() error {
	g.tick++

	switch g.battleState {
	case StateTurnP1:
		if updatePlayerTurn(g.players[0], p1Keys) {
			g.battleState = StateTurnP2
		}

	case StateTurnP2:
		if updatePlayerTurn(g.players[1], p2Keys) {
			g.battleState = StateResolve
			g.phaseTick = 0
			g.maxAmpPerPlayer = [2]float64{}
		}

	case StateResolve:
		g.phaseTick++
		// 各プレイヤーの波を個別に最大振幅計測
		t := float64(g.tick) / 60.0
		for i, p := range g.players {
			if len(p.fired) == 0 {
				continue
			}
			for x := 0; x < screenWidth; x += 4 {
				v := math.Abs(CompositeAmplitude(p.fired, float64(x), t))
				if v > g.maxAmpPerPlayer[i] {
					g.maxAmpPerPlayer[i] = v
				}
			}
		}
		if g.phaseTick >= resolveTicks {
			g.applyDamage()
			g.battleState = StateShowDamage
			g.phaseTick = 0
		}

	case StateShowDamage:
		g.phaseTick++
		if g.phaseTick >= damageTicks {
			if g.hp[0] <= 0 || g.hp[1] <= 0 {
				switch {
				case g.hp[0] <= 0 && g.hp[1] <= 0:
					g.winner = -1
				case g.hp[0] <= 0:
					g.winner = 1
				default:
					g.winner = 0
				}
				g.battleState = StateGameOver
			} else {
				// 次のラウンドへ
				g.players[0].fired = nil
				g.players[1].fired = nil
				g.battleState = StateTurnP1
			}
		}

	case StateGameOver:
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.init()
		}
	}

	return nil
}

// applyDamage は余剰振幅方式でダメージを計算・適用する。
//
// 各プレイヤーの実効振幅 = 個別最大振幅 × 波種相性倍率
// 実効振幅が大きい方が「勝ち」、差分が負けた側へのダメージになる。
// 差が閾値以下なら相殺でノーダメージ。
func (g *Game) applyDamage() {
	const cancelThreshold = 5.0 // px 以下の差は相殺とみなす

	p1Type := lastFiredType(g.players[0])
	p2Type := lastFiredType(g.players[1])

	// 実効振幅（相性倍率を乗算）
	eff1 := g.maxAmpPerPlayer[0] * WaveTypeMultiplier(p1Type, p2Type)
	eff2 := g.maxAmpPerPlayer[1] * WaveTypeMultiplier(p2Type, p1Type)

	excess := eff1 - eff2
	g.lastDamage = [2]int{0, 0}

	switch {
	case excess > cancelThreshold:
		// P1 の波が勝ち → P2 がダメージ
		g.lastDamage[1] = CalcDamage(excess, 1.0)
		g.hp[1] -= g.lastDamage[1]
	case excess < -cancelThreshold:
		// P2 の波が勝ち → P1 がダメージ
		g.lastDamage[0] = CalcDamage(-excess, 1.0)
		g.hp[0] -= g.lastDamage[0]
		// |excess| <= threshold: 相殺、ノーダメージ
	}

	if g.hp[0] < 0 {
		g.hp[0] = 0
	}
	if g.hp[1] < 0 {
		g.hp[1] = 0
	}
}

func lastFiredType(p *Player) WaveType {
	if len(p.fired) == 0 {
		return WaveSin
	}
	return p.fired[len(p.fired)-1].Type
}

// --- 描画ヘルパー ---

func hpColor(hp int) color.NRGBA {
	if hp > 40 {
		return color.NRGBA{R: 68, G: 255, B: 136, A: 255} // 緑
	}
	return color.NRGBA{R: 255, G: 68, B: 68, A: 255} // 赤
}

func drawHPBar(screen *ebiten.Image, x, y float32, hp, maxHP int) {
	const barW, barH float32 = 200, 14
	vector.FillRect(screen, x, y, barW, barH, color.NRGBA{R: 40, G: 40, B: 40, A: 255}, false)
	fill := barW * float32(hp) / float32(maxHP)
	if fill < 0 {
		fill = 0
	}
	vector.FillRect(screen, x, y, fill, barH, hpColor(hp), false)
	vector.StrokeRect(screen, x, y, barW, barH, 1, color.NRGBA{R: 160, G: 160, B: 160, A: 200}, false)
}

// Draw はフレームごとの描画処理。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 10, G: 10, B: 26, A: 255})

	t := float64(g.tick) / 60.0
	baseline := float64(waveBaseline)

	drawWaveRange := func(waves []Wave, col color.NRGBA, thick float32, xMin, xMax int) {
		for x := xMin + 1; x <= xMax; x++ {
			y0 := baseline - CompositeAmplitude(waves, float64(x-1), t)
			y1 := baseline - CompositeAmplitude(waves, float64(x), t)
			vector.StrokeLine(screen, float32(x-1), float32(y0), float32(x), float32(y1), thick, col, true)
		}
	}

	// HP バー（常時表示）
	drawHPBar(screen, 8, 8, g.hp[0], InitialHP)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("P1  %d HP", g.hp[0]), 8, 24)
	drawHPBar(screen, float32(screenWidth)-208, 8, g.hp[1], InitialHP)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("P2  %d HP", g.hp[1]), screenWidth-208, 24)

	// 基準線・中央区切り線
	vector.StrokeLine(screen, 0, float32(baseline), screenWidth, float32(baseline), 1,
		color.NRGBA{R: 50, G: 50, B: 70, A: 180}, false)
	vector.StrokeLine(screen, screenWidth/2, 44, screenWidth/2, float32(screenHeight)-44, 1,
		color.NRGBA{R: 80, G: 80, B: 80, A: 120}, false)

	switch g.battleState {
	case StateTurnP1:
		// P1 プレビュー（左半分・静止・薄い）
		p1prev := g.players[0].pending
		p1prev.Frequency = 0
		c := g.players[0].waveColor
		drawWaveRange([]Wave{p1prev}, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 80}, 1, 0, screenWidth/2)
		// ターン表示
		ebitenutil.DebugPrintAt(screen, ">>> P1 TURN  [Space] to FIRE <<<", screenWidth/2-120, 44)
		// パラメータ表示
		pw := g.players[0].pending
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type:%-8s Amp:%.0f Freq:%.2f Phase:%.2fπ",
			pw.Type, pw.Amplitude, pw.Frequency, pw.Phase/math.Pi), 8, screenHeight-36)
		ebitenutil.DebugPrintAt(screen, "[1-6]Type  [W/S]Amp  [Q/E]Freq  [Z/X]Phase", 8, screenHeight-24)

	case StateTurnP2:
		// P1 の発射済み波をアニメーション表示
		if len(g.players[0].fired) > 0 {
			c := g.players[0].waveColor
			drawWaveRange(g.players[0].fired, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 160}, 2, 0, screenWidth-1)
		}
		// P2 プレビュー（右半分・静止・薄い）
		p2prev := g.players[1].pending
		p2prev.Frequency = 0
		c := g.players[1].waveColor
		drawWaveRange([]Wave{p2prev}, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 80}, 1, screenWidth/2, screenWidth-1)
		// ターン表示
		ebitenutil.DebugPrintAt(screen, ">>> P2 TURN  [Enter] to FIRE <<<", screenWidth/2-120, 44)
		pw := g.players[1].pending
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type:%-8s Amp:%.0f Freq:%.2f Phase:%.2fπ",
			pw.Type, pw.Amplitude, pw.Frequency, pw.Phase/math.Pi), screenWidth/2+8, screenHeight-36)
		ebitenutil.DebugPrintAt(screen, "[Num1-6]Type  [↑↓]Amp  [,.]Freq  [N/M]Phase", screenWidth/2+8, screenHeight-24)

	case StateResolve:
		// 合成波アニメーション
		var allFired []Wave
		for _, p := range g.players {
			allFired = append(allFired, p.fired...)
		}
		drawWaveRange(allFired, color.NRGBA{R: 255, G: 255, B: 220, A: 255}, 3, 0, screenWidth-1)
		ebitenutil.DebugPrintAt(screen, "--- CLASH! ---", screenWidth/2-50, 44)

	case StateShowDamage:
		// 合成波（静止）＋ダメージ数値
		var allFired []Wave
		for _, p := range g.players {
			allFired = append(allFired, p.fired...)
		}
		drawWaveRange(allFired, color.NRGBA{R: 255, G: 255, B: 220, A: 180}, 2, 0, screenWidth-1)

		p1t := lastFiredType(g.players[0])
		p2t := lastFiredType(g.players[1])
		eff1 := g.maxAmpPerPlayer[0] * WaveTypeMultiplier(p1t, p2t)
		eff2 := g.maxAmpPerPlayer[1] * WaveTypeMultiplier(p2t, p1t)
		excess := eff1 - eff2
		const cancelThreshold = 5.0

		switch {
		case excess > cancelThreshold:
			ebitenutil.DebugPrintAt(screen, "P1 WINS THE CLASH!", screenWidth/2-60, 44)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("-%d HP", g.lastDamage[1]), screenWidth-80, screenHeight/2-8)
		case excess < -cancelThreshold:
			ebitenutil.DebugPrintAt(screen, "P2 WINS THE CLASH!", screenWidth/2-60, 44)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("-%d HP", g.lastDamage[0]), 40, screenHeight/2-8)
		default:
			ebitenutil.DebugPrintAt(screen, "--- CANCELLED! ---", screenWidth/2-60, 44)
		}
		// デバッグ: 実効振幅と余剰を表示
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("P1 amp:%.1f  P2 amp:%.1f  excess:%.1f", g.maxAmpPerPlayer[0], g.maxAmpPerPlayer[1], excess),
			screenWidth/2-150, screenHeight/2+24)
		// 相性表示
		m1 := WaveTypeMultiplier(p1t, p2t)
		m2 := WaveTypeMultiplier(p2t, p1t)
		if m1 != 1.0 || m2 != 1.0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Multiplier  P1: x%.1f  P2: x%.1f", m1, m2),
				screenWidth/2-100, screenHeight/2+36)
		}

	case StateGameOver:
		switch g.winner {
		case 0:
			ebitenutil.DebugPrintAt(screen, "=== P1 WINS! ===", screenWidth/2-60, screenHeight/2-10)
		case 1:
			ebitenutil.DebugPrintAt(screen, "=== P2 WINS! ===", screenWidth/2-60, screenHeight/2-10)
		default:
			ebitenutil.DebugPrintAt(screen, "=== DRAW! ===", screenWidth/2-45, screenHeight/2-10)
		}
		ebitenutil.DebugPrintAt(screen, "[R] Restart", screenWidth/2-40, screenHeight/2+10)
	}
}

// Layout は論理画面サイズを返す。
func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}
