package main

import (
	"math"
	"sort"
)

type Escape struct {
	Player *Player
}

const maxHitWhenTrapped int = 3

func (e *Escape) Play() Move {
	front := e.Player.FindShooterOnDirection(e.Player.GetDirection())
	// back := e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward())
	left := e.Player.FindShooterOnDirection(e.Player.GetDirection().Left())
	right := e.Player.FindShooterOnDirection(e.Player.GetDirection().Right())

	// TODO hindari escape ke arah orang lagi perang

	var paths []Path // list of possible path
	validAdjacent := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithDiagonalAdjacents(), WithEmptyAdjacent())
	if len(front) > 0 {
		var newAdjacent []Point
		for _, fp := range front {
			// iterate cuma sampai lokasi player. kalau udah ketemu, lgsg break
			for _, adj := range validAdjacent {
				canAttack := fp.CanHitPoint(adj)
				if !canAttack {
					newAdjacent = append(newAdjacent, adj)
				}
			}
		}
		validAdjacent = newAdjacent
	}
	if len(left) > 0 {
		var newAdjacent []Point
		for _, fp := range left {
			// iterate cuma sampai lokasi player. kalau udah ketemu, lgsg break
			for _, adj := range validAdjacent {
				if !fp.CanHitPoint(adj) {
					newAdjacent = append(newAdjacent, adj)
				}
			}
		}
		validAdjacent = newAdjacent
	}
	if len(right) > 0 {
		var newAdjacent []Point
		for _, fp := range right {
			// iterate cuma sampai lokasi player. kalau udah ketemu, lgsg break
			for _, adj := range validAdjacent {
				if !fp.CanHitPoint(adj) {
					newAdjacent = append(newAdjacent, adj)
				}
			}
		}
		validAdjacent = newAdjacent
	}
	for _, adj := range validAdjacent {
		aStar := NewAStar(e.Player.Game.Arena)
		path, err := aStar.SearchPath(e.Player.GetPosition(), adj)
		if err == ErrPathNotFound {
			continue
		}
		paths = append(paths, path)
	}

	if len(paths) == 0 {
		if len(front) > 0 {
			e.Player.trappedCount++
			if e.Player.trappedCount > maxHitWhenTrapped {
				e.Player.trappedCount = 0
				if len(left) > 0 {
					return TurnLeft
				}
				if len(right) > 0 {
					return TurnRight
				}
			}
			return Throw
		} else if len(left) > 0 {
			return TurnLeft
		} else if len(right) > 0 {
			return TurnRight
		} else {
			return e.Player.Walk()
		}
	}


	e.Player.trappedCount = 0
	sort.Sort(byPathLength(paths))

	// nextPt := paths[0][1]
	// foreach path,
	// 	ambil titik kedua (asumsi titik pertama adalah source)
	// 	calculate movement needed utk kesana


	minPathLength := math.MaxInt
	var shortestPaths []Path
	for _, path := range paths {
		if minPathLength > len(path) {
			minPathLength = len(path)
		}
	}

	for _, path := range paths {
		if len(path) <= minPathLength {
			shortestPaths = append(shortestPaths, path)
		}
	}

	requiredMoves := make([][]Move, len(shortestPaths))
	for idx, aPath := range shortestPaths {
		nextPt := aPath[1]
		moves, err := e.Player.MoveToAdjacent(nextPt)
		if err != nil {
			continue
		}
		requiredMoves[idx] = moves
	}

	if len(requiredMoves) == 0 {
		return e.Player.Walk()
	}

	mostEfficientMoves := requiredMoves[0]
	for _, sc := range requiredMoves {
		if len(sc) < len(mostEfficientMoves) {
			mostEfficientMoves = sc
		}
	}

	return mostEfficientMoves[0]
}

type byPathLength []Path

func (a byPathLength) Len() int           { return len(a) }
func (a byPathLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a byPathLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 2022-08-23 05:10:46.295 ICT2022/08/22 22:10:46 http: panic serving 169.254.1.1:41708: runtime error: index out of range [0] with length 0 goroutine 9 [running]: net/http.(*conn).serve.func1(0xc00010a140) /usr/local/go/src/net/http/server.go:1795 +0x139 panic(0x734380, 0xc00001c200) /usr/local/go/src/runtime/panic.go:679 +0x1b2 main.(*Escape).Play(0xc000010148, 0xc000010148, 0xffffffffffffffa7) /src/app/state_escape.go:74 +0xc03 main.(*NormalStrategy).Play(0xc000010138, 0x4, 0xc000014240) /src/app/strategy.go:22 +0x64 main.(*Player).Play(...) /src/app/player.go:40 main.(*Server).Play(0xc000084300, 0xc0000f8440, 0x31, 0xc00001c1e0, 0x2, 0x4, 0xc0000f4bd0, 0xc000065500, 0x3e76fd7eb008) /src/app/main.go:138 +0x217 main.Server.UpdateArena.func1(0x7cea00, 0xc0000ea0e0, 0xc0000eeb00) /src/app/main.go:125 +0x265 net/http.HandlerFunc.ServeHTTP(0xc000072c50, 0x7cea00, 0xc0000ea0e0, 0xc0000eeb00) /usr/local/go/src/net/http/server.go:2036 +0x44 github.com/gorilla/mux.(*Router).ServeHTTP(0xc0000ce000, 0x7cea00, 0xc0000ea0e0, 0xc0000ee700) /go/pkg/mod/github.com/gorilla/mux@v1.8.0/mux.go:210 +0xe2 net/http.serverHandler.ServeHTTP(0xc0000ea000, 0x7cea00, 0xc0000ea0e0, 0xc0000ee700) /usr/local/go/src/net/http/server.go:2831 +0xa4 net/http.(*conn).serve(0xc00010a140, 0x7cf080, 0xc0000582c0) /usr/local/go/src/net/http/server.go:1919 +0x875 created by net/http.(*Server).Serve /usr/local/go/src/net/http/server.go:2957 +0x384
