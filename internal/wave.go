package game

import "math"

// WaveType は波の種類を表す。
type WaveType int

const (
	WaveSin      WaveType = iota // sin(x)
	WaveCos                      // cos(x)
	WaveTan                      // tan(x)  ※クランプあり
	WaveSquare                   // sign(sin(x))  矩形波
	WaveSawtooth                 // 鋸歯波
	WaveTriangle                 // 三角波
)

// String は WaveType の名称を返す。
func (wt WaveType) String() string {
	return [...]string{"Sin", "Cos", "Tan", "Square", "Sawtooth", "Triangle"}[wt]
}

// Wave は単一の波を表す。
type Wave struct {
	Amplitude float64  // 振幅 A (0.1 〜 3.0)
	Frequency float64  // 周波数 freq [Hz]
	Phase     float64  // 位相 φ [rad]
	Type      WaveType // 波の種類
	Direction int      // 伝播方向: +1=左→右, -1=右→左
}

// ValueAt は時刻 t、位置 x における波の変位を返す。
//
//	y(x, t) = A · f(dir·k·x - ω·t + φ)
//	k = 2π / λ,  ω = 2π·freq,  dir = Direction (±1)
func (w Wave) ValueAt(x, t float64) float64 {
	const lambda = 200.0 // 波長 [px]
	const tanClamp = 3.0 // tan のクランプ係数（画面外抑制）
	dir := float64(w.Direction)
	if dir == 0 {
		dir = 1
	}
	k := 2 * math.Pi / lambda
	omega := 2 * math.Pi * w.Frequency
	arg := dir*k*x - omega*t + w.Phase
	switch w.Type {
	case WaveCos:
		return w.Amplitude * math.Cos(arg)
	case WaveTan:
		v := w.Amplitude * math.Tan(arg)
		return math.Max(-w.Amplitude*tanClamp, math.Min(w.Amplitude*tanClamp, v))
	case WaveSquare:
		// sign(sin(x))
		if math.Sin(arg) >= 0 {
			return w.Amplitude
		}
		return -w.Amplitude
	case WaveSawtooth:
		// (arg mod 2π)/π - 1  → [-1, 1)
		norm := math.Mod(arg, 2*math.Pi)
		if norm < 0 {
			norm += 2 * math.Pi
		}
		return w.Amplitude * (norm/math.Pi - 1)
	case WaveTriangle:
		// 2/π · arcsin(sin(x))
		return w.Amplitude * (2 / math.Pi) * math.Asin(math.Sin(arg))
	default: // WaveSin
		return w.Amplitude * math.Sin(arg)
	}
}

// CompositeAmplitude は複数の波を重ね合わせた振幅を返す。
func CompositeAmplitude(waves []Wave, x, t float64) float64 {
	var sum float64
	for _, w := range waves {
		sum += w.ValueAt(x, t)
	}
	return sum
}
