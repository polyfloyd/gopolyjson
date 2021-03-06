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

type Union struct {
	A, B Shape
}

type Shape interface {
	xxxShape()
}

func (Triangle) xxxShape() {}
func (Square) xxxShape()   {}
func (Polygon) xxxShape()  {}
func (Circle) xxxShape()   {}
func (Union) xxxShape()    {}

type Area struct {
	Color string
	Shape Shape `json:"shape"`
}

type Pattern struct {
	Size   int
	Shapes []Shape `json:"shapes"`
}

type NamedPattern struct {
	Sizes  map[string]int   `json:"sizes"`
	Shapes map[string]Shape `json:"named_shapes"`
}

type ShapeShifter struct {
	From, To Shape
	SkipMe   Shape `json:"-"`
	Err      string
}
