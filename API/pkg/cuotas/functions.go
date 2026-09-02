package cuotas

type genNum struct {
	top   int
	start int
}

func newGenNum(top int) *genNum {
	return &genNum{top: top, start: 0}
}

func (g *genNum) Next() int {
	if g.start < g.top {
		g.start++
	}
	return g.start
}
