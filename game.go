package main

import (
	"context"
)

type ArenaUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena struct {
		Dimensions []int                  `json:"dims"`
		State      map[string]PlayerState `json:"state"`
	} `json:"arena"`
}

type PlayerState struct {
	URL       string `json:"-"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
}

func (p PlayerState) GetDirection() Direction {
	switch p.Direction {
	case "N":
		return North
	case "W":
		return West
	case "E":
		return East
	default:
		return South
	}
}

type Mode string

func (m Mode) NeedLeaderboard() bool {
	return m == ZombieMode || m == AggressiveMode
}

const (
	NormalMode     Mode = "normal"
	BraveMode      Mode = "brave"
	// ZombieMode tries to attack player with lowest rank
	ZombieMode     Mode = "zombie"
	// AggressiveMode tries to climbing up the leaderboard
	AggressiveMode Mode = "aggressive"
)

type Game struct {
	Arena            Arena
	PlayerStateByURL map[string]PlayerState
	LeaderBoard      LeaderBoard
	Mode             Mode
}

const (
	defaultAttackRange int = 3
)

type GameOption func(*GameOptions)

type GameOptions struct {
	Mode Mode
}

func WithGameMode(m Mode) GameOption {
	return func(options *GameOptions) {
		options.Mode = m
	}
}

func NewGame(opts ...GameOption) Game {
	o := &GameOptions{
		Mode: NormalMode,
	}
	for _, opt := range opts {
		opt(o)
	}

	return Game{
		Mode: o.Mode,
	}
}

func (g *Game) UpdateArena(ctx context.Context, a ArenaUpdate) {
	ctx, span := tracer.Start(ctx, "Game.UpdateArena")
	defer span.End()

	width := a.Arena.Dimensions[0]
	height := a.Arena.Dimensions[1]
	arena := NewArena(width, height)

	g.LeaderBoard = []PlayerState{}
	for k, v := range a.Arena.State {
		v.URL = k
		arena.PutPlayer(v)
		g.LeaderBoard = append(g.LeaderBoard, v)
	}

	g.Arena = arena
	g.PlayerStateByURL = a.Arena.State
	if g.Mode.NeedLeaderboard() {
		g.UpdateLeaderBoard(ctx)
	}
}

func (g Game) Player(url string) *Player {
	pState := g.PlayerStateByURL[url]
	player := NewPlayerWithUrl(url, pState)
	player.Game = g
	// TODO set player strategy

	switch g.Mode {
	// case ZombieMode:
	// 	player.Strategy = NewChaseLowestRankStrategy()
	// case AggressiveMode:
		// rank := g.LeaderBoard.GetRank(*p)
		// if rank == 0 {
		// 	p.Strategy = NewNormalStrategy()
		// } else {
		// 	target := p.GetPlayerOnNextPodium(ctx)
		// 	p.Strategy = NewSafeChasing(target)
		// }
		// player.Strategy = NewNormalStrategy()
	case BraveMode:
		player.Strategy = NewBraveStrategy()
	default:
		player.Strategy = NewNormalStrategy()
	}
	return player
}

func (g Game) Update(player *Player) {
	updatedPlayer := g.PlayerStateByURL[player.Name]
	player.X = updatedPlayer.X
	player.Y = updatedPlayer.Y
	player.Direction = updatedPlayer.Direction
	player.WasHit = updatedPlayer.WasHit
	player.Score = updatedPlayer.Score
	player.Game = g
}

func (g *Game) UpdateLeaderBoard(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Game.UpdateLeaderBoard")
	defer span.End()

	g.LeaderBoard.Sort()
}

func (g Game) GetPlayerStateByPosition(p Point) (PlayerState, bool) {
	player := g.Arena.Grid[p.Y][p.X].Player
	if player == nil {
		return PlayerState{}, false
	}
	return *player, true
}

// GetPlayerByRank rank starts from 0 (highest rank)
func (g Game) GetPlayerByRank(rank int) *Player {
	ps := g.LeaderBoard.GetPlayerByRank(rank)
	if ps == nil {
		return nil
	}
	return g.GetPlayerByPosition(Point{ps.X, ps.Y})
}

func (g Game) GetPlayerByPosition(p Point) *Player {
	pState := g.Arena.Grid[p.Y][p.X].Player
	if pState == nil {
		return nil
	}
	player := NewPlayerWithUrl(pState.URL, *pState)
	player.Game = g
	return player
}

// ObstacleMap return map which denotes whether a cell is obstacle or not
func (g Game) ObstacleMap(ctx context.Context) [][]bool {
	ctx, span := tracer.Start(ctx, "Game.ObstacleMap")
	defer span.End()

	m := make([][]bool, g.Arena.Height)
	for i, _ := range m {
		row := make([]bool, g.Arena.Width)
		m[i] = row
	}

	for _, ps := range g.PlayerStateByURL {
		// m[ps.Y][ps.X] = true
		m[ps.Y][ps.X] = m[ps.Y][ps.X] || true

		if !ps.WasHit {
			continue
		}

		player := g.GetPlayerByPosition(Point{ps.X, ps.Y})
		if player == nil {
			continue
		}

		left := player.FindShooterOnDirection(ctx, player.GetDirection().Left())
		if left != nil {
			npt := left.GetPosition()
			for ctr := 0; ctr < defaultAttackRange; ctr++ {
				npt = npt.TranslateToDirection(1, left.GetDirection())
				if !g.Arena.IsValid(npt) {
					break
				}
				m[npt.Y][npt.X] = m[npt.Y][npt.X] || left.CanHitPoint(ctx, npt)
			}
		}

		front := player.FindShooterOnDirection(ctx, player.GetDirection())
		if front != nil {
			npt := front.GetPosition()
			for ctr := 0; ctr < defaultAttackRange; ctr++ {
				npt = npt.TranslateToDirection(1, front.GetDirection())
				if !g.Arena.IsValid(npt) {
					break
				}
				m[npt.Y][npt.X] = m[npt.Y][npt.X] || front.CanHitPoint(ctx, npt)
			}
		}

		back := player.FindShooterOnDirection(ctx, player.GetDirection().Backward())
		if back != nil {
			npt := back.GetPosition()
			for ctr := 0; ctr < defaultAttackRange; ctr++ {
				npt = npt.TranslateToDirection(1, back.GetDirection())
				if !g.Arena.IsValid(npt) {
					break
				}
				m[npt.Y][npt.X] = m[npt.Y][npt.X] || back.CanHitPoint(ctx, npt)
			}
		}

		right := player.FindShooterOnDirection(ctx, player.GetDirection().Right())
		if right != nil {
			npt := right.GetPosition()
			for ctr := 0; ctr < defaultAttackRange; ctr++ {
				npt = npt.TranslateToDirection(1, right.GetDirection())
				if !g.Arena.IsValid(npt) {
					break
				}
				m[npt.Y][npt.X] = m[npt.Y][npt.X] || right.CanHitPoint(ctx, npt)
			}
		}
	}

	return m
}
