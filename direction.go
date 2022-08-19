package main


type Direction struct {
	Name   string
	Degree int
}

func (d Direction) Left() Direction {
	degree := d.Degree - 90
	return DirectionMap[degree]
}

func (d Direction) Right() Direction {
	degree := d.Degree + 90
	return DirectionMap[degree]
}

func (d Direction) Backward() Direction {
	degree := d.Degree + 180
	return DirectionMap[degree]
}

var (
	North = Direction{"N", 270} //
	West  = Direction{"W", 180}
	South = Direction{"S", 90}
	East  = Direction{"E", 0}

	DirectionMap = map[int]Direction{
		-90: North,
		0:   East,
		90:  South,
		180: West,
		270: North,
		360: East,
		450: South,
		540: West,
	}
)


