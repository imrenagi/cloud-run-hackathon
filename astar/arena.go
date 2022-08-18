package astar

import (
	"fmt"
	"math"
	"sort"

	"github.com/rs/zerolog/log"
)

type Point struct {
	X, Y int
}

func (p Point) Equal(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
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
	F, G, H float64
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

func NewAStart() AStar {
	return AStar{distanceCalculator: &ManhattanDistance{}}
}

type AStar struct {
	distanceCalculator DistanceCalculator
}

func (as AStar) SearchPath(a Arena, src, dest Point) ([]Point, error) {

	// If the source is out of range
	if !a.IsValid(src) {
		log.Error().Msg("source is invalid")
		return nil, fmt.Errorf("source is invalid")
	}

	// If the destination is out of range
	if !a.IsValid(dest) {
		log.Error().Msg("destination is invalid")
		return nil, fmt.Errorf("destination is invalid")
	}

	// Either the source or the destination is blocked
	// TODO need to update this logic since the destination will be blocked because the player is the target
	if !a.Grid.IsUnblock(src) || !a.Grid.IsUnblock(dest) {
		log.Warn().Msg("Source or the destination is blocked")
		return nil, fmt.Errorf("source or destination is blocked")
	}

	// If the destination cell is the same as source cell
	if src.Equal(dest) {
		log.Info().Msg("We are already at the destination")
		return nil, nil
	}

	// Create a closed list and initialise it to false which
	// means that no cell has been included yet This closed
	// list is implemented as a boolean 2D array
	closedList := make([][]bool, a.Height)
	for y := range closedList {
		closedList[y] = make([]bool, a.Width)
	}

	// Declare a 2D array of structure to hold the details
	// of that cell
	cellDetails := make([][]Cell, a.Height)
	for y := range cellDetails {
		cellDetails[y] = make([]Cell, a.Width)
	}

	var i, j int

	for i := 0; i < a.Height; i++ {
		for j := 0; j < a.Width; j++ {
			cellDetails[i][j].F = math.MaxFloat64
			cellDetails[i][j].G = math.MaxFloat64
			cellDetails[i][j].H = math.MaxFloat64
			cellDetails[i][j].ParentY = -1
			cellDetails[i][j].ParentX = -1
		}
	}

	// Initialising the parameters of the starting node
	i = src.Y
	j = src.X
	cellDetails[i][j].F = 0.0
	cellDetails[i][j].G = 0.0
	cellDetails[i][j].H = 0.0
	cellDetails[i][j].ParentY = i
	cellDetails[i][j].ParentX = j

	/*
	   Create an open list having information as-
	   <f, <i, j>>
	   where f = g + h,
	   and i, j are the row and column index of that cell
	   Note that 0 <= i <= ROW-1 & 0 <= j <= COL-1
	   This open list is implemented as a set of pair of
	   pair.*/
	var openList []ppair
	openList = append(openList, ppair{F: 0, X: src.X, Y: src.Y})

	var foundDest bool

	for {
		if len(openList) == 0 {
			break
		}

		sort.Sort(byF(openList))
		p := openList[0]
		openList = openList[1:]


		i := p.Y
		j := p.X
		closedList[i][j] = true

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

		// To store the 'g', 'h' and 'f' of the 8 successors
		var gNew, hNew, fNew float64

		// Only process this cell if this is a valid one
		north := Point{X: j, Y: i - 1}
		if a.IsValid(north) {
			// If the destination cell is the same as the
			// current successor
			if a.IsDestination(north, dest) {
				cellDetails[north.Y][north.X].ParentY = i
				cellDetails[north.Y][north.X].ParentX = j
				log.Info().Msg("The destination cell is found")
				// tracePath(cellDetails, dest)
				foundDest = true
				break
			} else if !closedList[north.Y][north.X] && a.Grid.IsUnblock(north) {
				// If the successor is already on the closed
				// list or if it is blocked, then ignore it.
				// Else do the following
				// TODO calculate turn needed
				gNew = cellDetails[i][j].G + 1.0
				hNew = as.distanceCalculator.Distance(north, dest)
				fNew = gNew + hNew

				// If it isn’t on the open list, add it to
				// the open list. Make the current square
				// the parent of this square. Record the
				// f, g, and h costs of the square cell
				//                OR
				// If it is on the open list already, check
				// to see if this path to that square is
				// better, using 'f' cost as the measure.

				if cellDetails[north.Y][north.X].F == math.MaxFloat64 ||
				  cellDetails[north.Y][north.X].F > fNew {
					openList = append(openList, ppair{
						F: fNew,
						X: north.X,
						Y: north.Y,
					})

					cellDetails[north.Y][north.X].F = fNew
					cellDetails[north.Y][north.X].G = gNew
					cellDetails[north.Y][north.X].H = hNew
					cellDetails[north.Y][north.X].ParentY = i
					cellDetails[north.Y][north.X].ParentX = j
				}
			}
		}

		// Only process this cell if this is a valid one
		east := Point{X: j+1, Y: i}
		if a.IsValid(east) {
			if a.IsDestination(east, dest) {
				cellDetails[east.Y][east.X].ParentY = i
				cellDetails[east.Y][east.X].ParentX = j
				log.Info().Msg("The destination cell is found")
				// tracePath(cellDetails, dest)
				foundDest = true
				break
			} else if !closedList[east.Y][east.X] && a.Grid.IsUnblock(east) {
				// TODO calculate turn needed
				gNew = cellDetails[i][j].G + 1.0
				hNew = as.distanceCalculator.Distance(east, dest)
				fNew = gNew + hNew

				if cellDetails[east.Y][east.X].F == math.MaxFloat64 ||
				  cellDetails[east.Y][east.X].F > fNew {
					openList = append(openList, ppair{
						F: fNew,
						X: east.X,
						Y: east.Y,
					})
					cellDetails[east.Y][east.X].F = fNew
					cellDetails[east.Y][east.X].G = gNew
					cellDetails[east.Y][east.X].H = hNew
					cellDetails[east.Y][east.X].ParentY = i
					cellDetails[east.Y][east.X].ParentX = j
				}
			}
		}
	}

	if !foundDest {
		return nil, fmt.Errorf("path not found")
	}

	return as.tracePath(cellDetails, dest), nil
}

func (as AStar) tracePath(cellDetails [][]Cell, dest Point) []Point{
	var finalPath []Point
	row := dest.Y
	col := dest.X
	var path []Point //stack
	for {
		if !(!(cellDetails[row][col].ParentY == row && cellDetails[row][col].ParentX == col)) {
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
		fmt.Printf("-> (%d, %d)", p.X, p.Y)
		finalPath = append(finalPath, p)
	}
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
	return math.Abs(float64(p1.X-p2.X) + math.Abs(float64(p1.Y-p2.Y)))
}

type EuclideanDistance struct {
}

func (e EuclideanDistance) Distance(p1, p2 Point) float64 {
	return math.Sqrt(float64(p1.X-p2.X)*float64(p1.X-p2.X) + float64(p1.Y-p2.Y)*float64(p1.Y-p2.Y))
}
