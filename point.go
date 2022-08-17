package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p Point) Equal(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p Point) MoveWithDirection(distance int, direction Direction) Point {
	dir := float64(direction.Degree) * math.Pi / 180
	return Point{
		X: p.X + distance*int(math.Round(math.Sin(dir))),
		Y: p.Y + distance*int(math.Round(math.Cos(dir))),
	}
}