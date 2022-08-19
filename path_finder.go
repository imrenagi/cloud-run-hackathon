package main

import "strings"

type Path []Point

func (p Path) String() string {
	var ps []string
	for _, pt := range p {
		ps = append(ps, pt.String())
	}
	return strings.Join(ps, "->")
}

type PathFinder interface {
	SearchPath(src, dest Point) (Path, error)
}
