package testdata

type Triangle struct {
	P0 [2]int
	P1 [2]int
	P2 [2]int
}

type Square struct {
	TopLeft       [2]int
	Width, Height int
}

type (
	Polygon struct {
		Vertices [][2]int
	}
	Circle struct {
		Center [2]int
		Radius int
	}
)

type Shape interface {
	xxxShape()
}

func (Triangle) xxxShape() {}
func (Square) xxxShape()   {}
func (Polygon) xxxShape()  {}
func (Circle) xxxShape()   {}

type Area struct {
	Color string
	Shape Shape `json:"shape"`
}

type Pattern struct {
	Size   int
	Shapes []Shape `json:"shapes"`
}

type ShapeShifter struct {
	From, To Shape
	SkipMe   Shape `json:"-"`
}
