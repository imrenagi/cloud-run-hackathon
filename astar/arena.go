package astar

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

type Point struct {
	X, Y int
}

func (p Point) Equal(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

type Grid [][]Cell

// A Utility Function to check whether the given cell is
// blocked or not
func (g Grid) IsUnblock(p Point) bool {
	return g[p.Y][p.X].Player == nil
}

type Arena struct {
	Width  int // x start from top left to the right
	Height int // y start from top left to the bottom
	Grid   Grid
}

type Player struct {
	Name string
}

type Cell struct {
	Player  *Player
	ParentX int
	ParentY int
	F float64
	G, H float64
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

type Option func(options *Options)

type Options struct {
	DistanceCalculator DistanceCalculator
}

func NewAStar(a Arena, opts ...Option) AStar {

	options := &Options{
		DistanceCalculator: &ManhattanDistance{},
	}
	for _, o := range opts {
		o(options)
	}

	closedList := make([][]bool, a.Height)
	for y := range closedList {
		closedList[y] = make([]bool, a.Width)
	}

	cellDetails := make([][]Cell, a.Height)
	for y := range cellDetails {
		cellDetails[y] = make([]Cell, a.Width)
	}

	return AStar{
		arena:              a,
		distanceCalculator: options.DistanceCalculator,
		closedList:         closedList,
		cellDetails:        cellDetails,
	}
}

type AStar struct {
	distanceCalculator DistanceCalculator
	arena              Arena
	// Create a closed list and initialise it to false which
	// means that no cell has been included yet This closed
	// list is implemented as a boolean 2D array
	closedList [][]bool
	// Declare a 2D array of structure to hold the details
	// of that cell
	cellDetails [][]Cell
	/*
	   Create an open list having information as-
	   <f, <i, j>>
	   where f = g + h,
	   and i, j are the row and column index of that cell
	   Note that 0 <= i <= ROW-1 & 0 <= j <= COL-1
	   This open list is implemented as a set of pair of
	   pair.*/
	openList []ppair
}

func (as *AStar) SearchPath(src, dest Point) ([]Point, error) {

	// If the source is out of range
	if !as.arena.IsValid(src) {
		log.Error().Msg("source is invalid")
		return nil, fmt.Errorf("source is invalid")
	}

	// If the destination is out of range
	if !as.arena.IsValid(dest) {
		log.Error().Msg("destination is invalid")
		return nil, fmt.Errorf("destination is invalid")
	}

	// Either the source or the destination is blocked
	// TODO need to update this logic since the destination will be blocked because the player is the target
	if !as.arena.Grid.IsUnblock(src) || !as.arena.Grid.IsUnblock(dest) {
		log.Warn().Msg("Source or the destination is blocked")
		return nil, fmt.Errorf("source or destination is blocked")
	}

	// If the destination cell is the same as source cell
	if src.Equal(dest) {
		log.Info().Msg("We are already at the destination")
		return nil, nil
	}

	for i := 0; i < as.arena.Height; i++ {
		for j := 0; j < as.arena.Width; j++ {
			as.cellDetails[i][j].F = math.MaxFloat64
			as.cellDetails[i][j].G = math.MaxFloat64
			as.cellDetails[i][j].H = math.MaxFloat64
			as.cellDetails[i][j].ParentY = -1
			as.cellDetails[i][j].ParentX = -1
		}
	}

	// Initialising the parameters of the starting node
	as.cellDetails[src.Y][src.X].F = 0.0
	as.cellDetails[src.Y][src.X].G = 0.0
	as.cellDetails[src.Y][src.X].H = 0.0
	as.cellDetails[src.Y][src.X].ParentY = src.Y
	as.cellDetails[src.Y][src.X].ParentX = src.X

	as.openList = append(as.openList, ppair{F: 0, X: src.X, Y: src.Y})

	var foundDest bool

	for {
		if len(as.openList) == 0 {
			break
		}
		sort.Sort(byF(as.openList))
		currNode := as.openList[0]
		as.openList = as.openList[1:]
		as.closedList[currNode.Y][currNode.X] = true
		/*
		   Generating all the 8 successor of this cell

		       N.W   N   N.E
		         \   |   /
		          \  |  /
		       W----Cell----E
		            / | \
		          /   |  \
		       S.W    S   S.E

		   Cell-->Popped Cell (i, j)
		   N -->  North       (i-1, j)
		   S -->  South       (i+1, j)
		   E -->  East        (i, j+1)
		   W -->  West           (i, j-1)
		   N.E--> North-East  (i-1, j+1)
		   N.W--> North-West  (i-1, j-1)
		   S.E--> South-East  (i+1, j+1)
		   S.W--> South-West  (i+1, j-1)*/

		// Only process this cell if this is a valid one
		north := Point{X: currNode.X, Y: currNode.Y - 1}
		foundDest = as.checkSuccessor(currNode, north, dest)
		if foundDest {
			break
		}

		south := Point{X: currNode.X, Y: currNode.Y + 1}
		foundDest = as.checkSuccessor(currNode, south, dest)
		if foundDest {
			break
		}

		east := Point{X: currNode.X + 1, Y: currNode.Y}
		foundDest = as.checkSuccessor(currNode, east, dest)
		if foundDest {
			break
		}

		west := Point{X: currNode.X - 1, Y: currNode.Y}
		foundDest = as.checkSuccessor(currNode, west, dest)
		if foundDest {
			break
		}

		// because we only move to horizontal and vertical, we disabled this
		// TODO create function to get allowed direction. moving backward

		// northEast := Point{X: currNode.X + 1, Y: currNode.Y - 1}
		// foundDest = as.checkSuccessor(currNode, northEast, dest)
		// if foundDest {
		// 	break
		// }
		//
		// northWest := Point{X: currNode.X - 1, Y: currNode.Y - 1}
		// foundDest = as.checkSuccessor(currNode, northWest, dest)
		// if foundDest {
		// 	break
		// }
		//
		// southEast := Point{X: currNode.X + 1, Y: currNode.Y + 1}
		// foundDest = as.checkSuccessor(currNode, southEast, dest)
		// if foundDest {
		// 	break
		// }
		//
		// southWest := Point{X: currNode.X - 1, Y: currNode.Y + 1}
		// foundDest = as.checkSuccessor(currNode, southWest, dest)
		// if foundDest {
		// 	break
		// }
	}

	if !foundDest {
		return nil, fmt.Errorf("path not found")
	}

	return as.tracePath(as.cellDetails, dest), nil
}

func (as *AStar) checkSuccessor(currNode ppair, successor Point, dest Point) bool {

	var gNew, hNew, fNew float64
	// Only process this cell if this is a valid one
	if as.arena.IsValid(successor) {
		if as.arena.IsDestination(successor, dest) {
			as.cellDetails[successor.Y][successor.X].ParentY = currNode.Y
			as.cellDetails[successor.Y][successor.X].ParentX = currNode.X
			log.Info().Msg("The destination cell is found")
			return true
		} else if !as.closedList[successor.Y][successor.X] && as.arena.Grid.IsUnblock(successor) {
			// TODO calculate turn needed
			gNew = as.cellDetails[currNode.Y][currNode.X].G + 1.0
			hNew = as.distanceCalculator.Distance(successor, dest)
			fNew = gNew + hNew

			if as.cellDetails[successor.Y][successor.X].F == math.MaxFloat64 ||
			  as.cellDetails[successor.Y][successor.X].F > fNew {
				as.openList = append(as.openList, ppair{
					F: fNew,
					X: successor.X,
					Y: successor.Y,
				})
				as.cellDetails[successor.Y][successor.X].F = fNew
				as.cellDetails[successor.Y][successor.X].G = gNew
				as.cellDetails[successor.Y][successor.X].H = hNew
				as.cellDetails[successor.Y][successor.X].ParentY = currNode.Y
				as.cellDetails[successor.Y][successor.X].ParentX = currNode.X
			}
		}
	}
	return false
}

type Path []Point

func (p Path) String() string {
	var ps []string
	for _, pt := range p {
		ps = append(ps, pt.String())
	}
	return strings.Join(ps, "->")
}

func (as AStar) tracePath(cellDetails [][]Cell, dest Point) []Point {
	var finalPath Path
	row := dest.Y
	col := dest.X
	var path []Point // stack
	for {
		if cellDetails[row][col].ParentY == row && cellDetails[row][col].ParentX == col {
			break
		}
		path = append(path, Point{
			X: col,
			Y: row,
		})
		tempRow := cellDetails[row][col].ParentY
		tempCol := cellDetails[row][col].ParentX
		row = tempRow
		col = tempCol
	}

	path = append(path, Point{
		X: col,
		Y: row,
	})
	for {
		if len(path) == 0 {
			break
		}
		p := path[len(path)-1]
		path = path[:len(path)-1]
		finalPath = append(finalPath, p)
	}

	log.Info().Stringer("path", finalPath).Msg("path found")
	return finalPath
}

type ppair struct {
	F    float64
	X, Y int
}

// ByAge implements sort.Interface based on the Age field.
type byF []ppair

func (a byF) Len() int           { return len(a) }
func (a byF) Less(i, j int) bool { return a[i].F < a[j].F }
func (a byF) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ===== distance ====

type DistanceCalculator interface {
	Distance(src, dest Point) float64
}

type ManhattanDistance struct {
}

func (m ManhattanDistance) Distance(p1, p2 Point) float64 {
	return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}

type EuclideanDistance struct {
}

func (e EuclideanDistance) Distance(p1, p2 Point) float64 {
	return math.Sqrt(float64(p1.X-p2.X)*float64(p1.X-p2.X) + float64(p1.Y-p2.Y)*float64(p1.Y-p2.Y))
}
