package main

type Grid [][]Cell // TODO should this be [][]Cell or [][]Player

// A Utility Function to check whether the given cell is
// blocked or not
func (g Grid) IsUnblock(p Point) bool {
	return g[p.Y][p.X].Player == nil
}

type Cell struct {
	// TODO Should cell know about its location?
	Player  *PlayerState
}

func NewArena(w, h int) Arena {
	grid := make([][]Cell, h)
	for y := range grid {
		grid[y] = make([]Cell, w)
		for x := range grid[y] {
			grid[y][x] = Cell{}
		}
	}
	return Arena{
		Width:  w,
		Height: h,
		Grid:   grid,
	}
}

type Arena struct {
	Width  int // x start from top left to the right
	Height int // y start from top left to the bottom
	Grid   Grid
}

func (a *Arena) PutPlayer(p PlayerState) {
	a.Grid[p.Y][p.X].Player = &p
}

// Traverse with BFS
func (a Arena) Traverse(start Point) []Point {
	var traversedNode []Point

	visited := make([]bool, a.Width*a.Height)
	queue := make([]Point, 0)
	visited[start.Y*a.Width+start.X] = true
	queue = append(queue, start)

	for {
		if len(queue) == 0 {
			break
		}
		pt := queue[0]
		traversedNode = append(traversedNode, pt)
		queue = queue[1:]
		adjancentNodes := a.GetAdjacent(pt)
		for _, n := range adjancentNodes {
			if !visited[n.Y*a.Width+n.X] {
				visited[n.Y*a.Width+n.X] = true
				queue = append(queue, n)
			}
		}
	}

	return traversedNode
}

type AdjacentOption func (*AdjacentOptions)

type AdjacentOptions struct {
	IncludeDiagonal bool
}

func WithDiagonalAdjacents() AdjacentOption {
	return func(options *AdjacentOptions) {
		options.IncludeDiagonal = true
	}
}

func (a Arena) GetAdjacent(p Point, opts ...AdjacentOption) []Point {
	options := &AdjacentOptions{}
	for _, o := range opts {
		o(options)
	}
	var adj []Point
	iterator := [3]int{-1, 0, 1}
	for _, i := range iterator {
		for _, j := range iterator {
			if !options.IncludeDiagonal {
				if i * j == -1 || i*j == 1 {
					continue
				}
			}
			if i != 0 || j != 0 {
				p := Point{X: p.X + i, Y: p.Y + j}
				if p.IsInArena(a) {
					adj = append(adj, p)
				}
			}
		}
	}
	return adj
}

// A Utility Function to check whether given cell (row, col)
// is a valid cell or not.
func (a Arena) IsValid(p Point) bool {
	return p.Y >= 0 && p.Y < a.Height && p.X >= 0 && p.X < a.Width
}

// A Utility Function to check whether destination cell has
// been reached or not
func (a Arena) IsDestination(p, dest Point) bool {
	return p.Equal(dest)
}
