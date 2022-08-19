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

// TranslateAngle move point by using defined distance and direction
func (p Point) TranslateAngle(distance int, direction Direction) Point {
	dir := float64(direction.Degree) * math.Pi / 180
	return Point{
		X: p.X + distance*int(math.Round(math.Sin(dir))),
		Y: p.Y + distance*int(math.Round(math.Cos(dir))),
	}
}

// Translate moves p in x and y axis
func (p Point) Translate(x, y int) Point {
	return Point{X: p.X + x, Y: p.Y + y}
}

func (p Point) IsInArena(a Arena) bool {
	if p.X > a.Width-1 || p.X < 0 {
		return false
	}
	if p.Y > a.Height-1 || p.Y < 0 {
		return false
	}
	return true
}
