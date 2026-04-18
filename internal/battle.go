package game

const (
	InitialHP   = 100
	DamageCoeff = 0.4 // maxAmp[px] × 0.4 → HP ダメージ
)

// WaveTypeMultiplier は攻撃側の波種と防御側の波種の相性倍率を返す。
// 三すくみ: Sin→(有利)→Tan→(有利)→Square→(有利)→Sin
func WaveTypeMultiplier(attacker, defender WaveType) float64 {
	switch attacker {
	case WaveSin:
		switch defender {
		case WaveTan:
			return 1.5
		case WaveSquare:
			return 0.7
		}
	case WaveTan:
		switch defender {
		case WaveSquare:
			return 1.5
		case WaveSin:
			return 0.7
		}
	case WaveSquare:
		switch defender {
		case WaveSin:
			return 1.5
		case WaveTan:
			return 0.7
		}
	}
	return 1.0
}

// CalcDamage は合成波の最大振幅から HP ダメージを算出する。
//
//	damage = maxAmp × DamageCoeff × multiplier  (最低 1)
func CalcDamage(maxAmp, multiplier float64) int {
	d := int(maxAmp * DamageCoeff * multiplier)
	if d < 1 {
		d = 1
	}
	return d
}
