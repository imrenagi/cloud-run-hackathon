package main


type Direction struct {
	Name   string
	Degree int
}

func (d Direction) Left() Direction {
	degree := d.Degree + 90
	return DirectionMap[degree]
}

func (d Direction) Right() Direction {
	degree := d.Degree - 90
	return DirectionMap[degree]
}

func (d Direction) Backward() Direction {
	degree := d.Degree + 180
	return DirectionMap[degree]
}

var (
	North = Direction{"N", 180} //
	West  = Direction{"W", 270}
	South = Direction{"S", 0}
	East  = Direction{"E", 90}

	DirectionMap = map[int]Direction{
		-90: West,
		0:   South,
		90:  East,
		180: North,
		270: West,
		360: South,
		450: East,
		540: North,
	}
)
