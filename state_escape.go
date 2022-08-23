package main

import (
	"sort"

	"github.com/rs/zerolog/log"
)

type Escape struct {
	// TODO should pass pointer
	Player *Player
}

const maxHitWhenTrapped int = 3

func (e *Escape) Play() Move {
	front := e.Player.FindShooterOnDirection(e.Player.GetDirection())
	// back := e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward())
	left := e.Player.FindShooterOnDirection(e.Player.GetDirection().Left())
	right := e.Player.FindShooterOnDirection(e.Player.GetDirection().Right())


	var paths []Path // list of possible path
	validAdjacent := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithDiagonalAdjacents(), WithEmptyAdjacent())
	if len(front) > 0 {
		var newAdjacent []Point
		for _, fp := range front {
			// iterate cuma sampai lokasi player. kalau udah ketemu, lgsg break
			for _, adj := range validAdjacent {
				if !fp.CanAttack(adj) {
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
				if !fp.CanAttack(adj) {
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
				if !fp.CanAttack(adj) {
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
		log.Info().
			Int("trappedCount", e.Player.trappedCount).
			Bool("wasHit", e.Player.WasHit).
			Int("score", e.Player.Score).
			Msgf("trapped")

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

	// foreach path,
	// 	ambil titik kedua (asumsi titik pertama adalah source)
	// 	calculate movement needed utk kesana
	requiredMoves := make([][]Move, len(paths))

	// TODO ambil dulu path terpendek
	for idx, aPath := range paths {
		nextPt := aPath[1]
		moves, err := e.Player.MoveNeededToReachAdjacent(nextPt)
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

	//
	// // TODO cari adjacent dengan movement paling minimal
	// // isOnAttackRange := make([]bool, len(emptyAdjacents))
	// requiredMoves := make([][]Move, len(emptyAdjacents))
	// for idx, adj := range emptyAdjacents {
	// 	decisions, err := e.Player.MoveNeededToReachAdjacent(adj)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	requiredMoves[idx] = decisions
	// }
	//
	// //untuk setiap adjacent, run a star, cari path terpendek dan f value terkecil
	// //kalau path not found, tembak 3 kali.
	//
	// // TODO: cari 8 adjacent. Cari adjacent yg aman dari tembakan front, back, left, right.
	// // cari shortest path, execute.
	//
	// // TODO bug: kalau terpojok, tapi dari jarak tertentu. jadinya masih ada adjacent yg bisa
	// // tapi malah stuck puter2
	//
	// // TODO bug: kalau hadap2an, dia malah ngejar. harusnya kabur. ini karena adjacent di depan
	// // available dan satisfy the minimum step juga lebih kecil.
	//
	// // TODO check whether adjacent is on opponents attack range?
	// // kalau gak ada ambil yg ada aja.
	//
	// if len(requiredMoves) == 0 {
	// 	return e.Player.Walk()
	// }
	//
	// mostEfficientMoves := requiredMoves[0]
	// for _, sc := range requiredMoves {
	// 	if len(sc) < len(mostEfficientMoves) {
	// 		mostEfficientMoves = sc
	// 	}
	// }
	//
	// return mostEfficientMoves[0]
	// return e.Player.Walk()
}

type byPathLength []Path

func (a byPathLength) Len() int           { return len(a) }
func (a byPathLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a byPathLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 2022-08-23 05:10:46.295 ICT2022/08/22 22:10:46 http: panic serving 169.254.1.1:41708: runtime error: index out of range [0] with length 0 goroutine 9 [running]: net/http.(*conn).serve.func1(0xc00010a140) /usr/local/go/src/net/http/server.go:1795 +0x139 panic(0x734380, 0xc00001c200) /usr/local/go/src/runtime/panic.go:679 +0x1b2 main.(*Escape).Play(0xc000010148, 0xc000010148, 0xffffffffffffffa7) /src/app/state_escape.go:74 +0xc03 main.(*NormalStrategy).Play(0xc000010138, 0x4, 0xc000014240) /src/app/strategy.go:22 +0x64 main.(*Player).Play(...) /src/app/player.go:40 main.(*Server).Play(0xc000084300, 0xc0000f8440, 0x31, 0xc00001c1e0, 0x2, 0x4, 0xc0000f4bd0, 0xc000065500, 0x3e76fd7eb008) /src/app/main.go:138 +0x217 main.Server.UpdateArena.func1(0x7cea00, 0xc0000ea0e0, 0xc0000eeb00) /src/app/main.go:125 +0x265 net/http.HandlerFunc.ServeHTTP(0xc000072c50, 0x7cea00, 0xc0000ea0e0, 0xc0000eeb00) /usr/local/go/src/net/http/server.go:2036 +0x44 github.com/gorilla/mux.(*Router).ServeHTTP(0xc0000ce000, 0x7cea00, 0xc0000ea0e0, 0xc0000ee700) /go/pkg/mod/github.com/gorilla/mux@v1.8.0/mux.go:210 +0xe2 net/http.serverHandler.ServeHTTP(0xc0000ea000, 0x7cea00, 0xc0000ea0e0, 0xc0000ee700) /usr/local/go/src/net/http/server.go:2831 +0xa4 net/http.(*conn).serve(0xc00010a140, 0x7cf080, 0xc0000582c0) /usr/local/go/src/net/http/server.go:1919 +0x875 created by net/http.(*Server).Serve /usr/local/go/src/net/http/server.go:2957 +0x384
