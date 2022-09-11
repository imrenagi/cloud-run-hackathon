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
	NormalMode Mode = "normal"
	BraveMode  Mode = "brave"
	// ZombieMode tries to attack player with lowest rank
	ZombieMode Mode = "zombie"
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

	switch g.Mode {
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
