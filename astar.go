package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/rs/zerolog/log"
)

type Option func(options *Options)

type Options struct {
	DistanceCalculator DistanceCalculator
}

type cellDetail struct {
	ParentX int
	ParentY int
	F       float64
	G, H    float64
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

	cellDetails := make([][]cellDetail, a.Height)
	for y := range cellDetails {
		cellDetails[y] = make([]cellDetail, a.Width)
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
	cellDetails [][]cellDetail
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
	// if !as.arena.Grid.IsUnblock(src) || !as.arena.Grid.IsUnblock(dest) {
	// 	log.Warn().Msg("Source or the destination is blocked")
	// 	return nil, fmt.Errorf("source or destination is blocked")
	// }

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

	player := as.arena.GetPlayer(src)

	as.openList = append(as.openList, ppair{F: 0, X: src.X, Y: src.Y, Direction: player.GetDirection()})

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
	}

	if !foundDest {
		return nil, ErrPathNotFound
	}

	return as.tracePath(as.cellDetails, dest), nil
}

var ErrPathNotFound = fmt.Errorf("path not found")

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
			step := currNode.requiredRotation(successor)

			gNew = as.cellDetails[currNode.Y][currNode.X].G + 1.0 + float64(step)
			hNew = as.distanceCalculator.Distance(successor, dest)
			fNew = gNew + hNew

			if as.cellDetails[successor.Y][successor.X].F == math.MaxFloat64 ||
			  as.cellDetails[successor.Y][successor.X].F > fNew {
				as.openList = append(as.openList, ppair{
					F: fNew,
					X: successor.X,
					Y: successor.Y,
					Direction: currNode.Direction,
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

func (as AStar) tracePath(cellDetails [][]cellDetail, dest Point) []Point {
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

	return finalPath
}

type ppair struct {
	F         float64
	X, Y      int
	Direction Direction
}

func (p *ppair) rotateCounterClockwise() {
	p.Direction = p.Direction.Left()
}

func (p *ppair) rotateClockwise() {
	p.Direction = p.Direction.Right()
}

func (p *ppair) setDirection(d Direction) {
	p.Direction = d
}

// requiredRotation returns the number of rotation need to be performed to head toward the pt
func (p *ppair) requiredRotation(pt Point) int {
	myPt := Point{X: p.X, Y: p.Y}
	const distance = 1
	var cCount, ccCount int // clockwise and counter clockwise counter
	initialDirection := p.Direction
	for i := 0; i<4; i++ {
		ptInFront := myPt.TranslateToDirection(distance, p.Direction)
		if ptInFront.Equal(pt) {
			break
		}
		p.rotateCounterClockwise()
		ccCount++
	}
	p.setDirection(initialDirection)
	for i := 0; i<4; i++ {
		ptInFront := myPt.TranslateToDirection(distance, p.Direction)
		if ptInFront.Equal(pt) {
			break
		}
		p.rotateClockwise()
		cCount++
	}
	p.setDirection(initialDirection)

	minRotationCount := ccCount
	for i := 0; i<ccCount; i++ {
		p.rotateCounterClockwise()
	}
	if minRotationCount > cCount {
		minRotationCount = cCount
		p.setDirection(initialDirection)
		for i := 0; i<cCount; i++ {
			p.rotateClockwise()
		}
	}
	return minRotationCount
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
