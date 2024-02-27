package orange

import "math"

var (
	Oranges = []*Orange{
		NewOrange(100, 0),
		NewOrange(110, -100),
		NewDekopon(100, -1000),
		NewHassaku(110, -2000),
		NewHyuganatu(120, -3000),
		NewLemon(150, -4000),
		NewMikan(170, -5000),
	}
)

const (
	gravity  = 0.98
	friction = 0.98
	bounce   = 0.3
	spring   = 0.4
)

type Calc struct {
	World World
	Score int
}

func (c *Calc) Oranges(oranges []*Orange) []*Orange {
	oranges = c.combine(oranges)
	c.move(oranges)
	c.Collision(oranges)
	c.screenWrap(oranges)

	return oranges
}

func (c *Calc) move(o []*Orange) {
	for _, orange := range Oranges {
		orange.VX *= friction
		orange.VY *= friction
		orange.VY += gravity
		orange.X += orange.VX
		orange.Y += orange.VY
	}
}

func (c *Calc) screenWrap(oranges []*Orange) {
	for _, orange := range oranges {
		// X軸方向にはみ出す
		if orange.X-orange.Radius < 0 {
			orange.X = orange.Radius
			orange.VX *= -bounce
		} else if c.World.Width < orange.X+orange.Radius {
			orange.X = c.World.Width - orange.Radius
			orange.VX *= -bounce
		}

		// Y軸方向にはみ出す
		if orange.Y < 0 {
			// なにもしない
		} else if c.World.Height < orange.Y+orange.Radius {
			orange.Y = c.World.Height - orange.Radius
			orange.VY *= -bounce
		}
	}
}

func (c *Calc) Collision(oranges []*Orange) {
	for i, orange := range oranges {
		for j := i + 1; j < len(oranges); j++ {
			g := oranges[j]
			dx := g.X - orange.X
			dy := g.Y - orange.Y
			d := math.Sqrt(dx*dx + dy*dy)
			minD := orange.Radius + g.Radius
			if d < minD {
				angle := math.Atan2(dy, dx)
				tx := orange.X + math.Cos(angle)*minD
				ty := orange.Y + math.Sin(angle)*minD
				ax := (tx - g.X) * spring
				ay := (ty - g.Y) * spring
				orange.VX -= ax
				orange.VY -= ay
				g.VX += ax
				g.VY += ay

				orange.X = orange.X - math.Cos(angle)*(minD-d)/2
				orange.Y = orange.Y - math.Sin(angle)*(minD-d)/2
				g.X = g.X + math.Cos(angle)*(minD-d)/2
				g.Y = g.Y + math.Sin(angle)*(minD-d)/2
			}
		}
	}
}

func (c *Calc) combine(oranges []*Orange) []*Orange {
	newOranges := make([]*Orange, 0)
	for i, orange := range oranges {
		for j := i + 1; j < len(oranges); j++ {
			g := oranges[j]
			if orange.Remove || g.Remove {
				continue
			}

			dx := g.X - orange.X
			dy := g.Y - orange.Y
			d := math.Sqrt(dx*dx + dy*dy)
			minD := orange.Radius + g.Radius
			if int64(d) <= int64(minD) && orange.Type == g.Type {
				orange.Remove = true
				g.Remove = true
				var next *Orange
				switch orange.Type {
				case ORANGE:
					next = NewHassaku((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 10
				case HASSAKU:
					next = NewDekopon((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 20
				case DEKOPON:
					next = NewHyuganatu((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 30
				case HYUGANATU:
					next = NewLemon((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 40
				case LEMON:
					next = NewMikan((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 50
				case MIKAN:
					next = NewMikan((orange.X+g.X)/2, (orange.Y+g.Y)/2)
					c.Score += 60
				}
				if next != nil {
					newOranges = append(newOranges, next)
				}
			}
		}
	}

	for _, orange := range oranges {
		if !orange.Remove {
			newOranges = append(newOranges, orange)
		}
	}

	return newOranges
}
