package orange

const (
	ORANGE = iota
	HASSAKU
	DEKOPON
	HYUGANATU
	SETOKA
	LEMON
	MIKAN
)

type Orange struct {
	X      float64
	Y      float64
	VX     float64
	VY     float64
	Radius float64
	Type   int
	Remove bool
}

func NewOrange(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 20,
		Type:   ORANGE,
	}
}

func NewHassaku(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 35,
		Type:   HASSAKU,
	}
}

func NewDekopon(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 50,
		Type:   DEKOPON,
	}
}

func NewHyuganatu(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 65,
		Type:   HYUGANATU,
	}
}

func NewLemon(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 80,
		Type:   LEMON,
	}
}

func NewMikan(x, y float64) *Orange {
	return &Orange{
		X:      x,
		Y:      y,
		VX:     0,
		VY:     0,
		Radius: 20,
		Type:   MIKAN,
	}
}
