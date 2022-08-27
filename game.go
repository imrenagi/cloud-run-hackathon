package main

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



type Game struct {
	Arena            Arena
	PlayerStateByURL map[string]PlayerState
	LeaderBoard      LeaderBoard
}

func NewGame() Game {
	return Game{
	}
}

func (g *Game) UpdateArena(a ArenaUpdate) {
	width := a.Arena.Dimensions[0]
	height := a.Arena.Dimensions[1]
	arena := NewArena(width, height)

	g.LeaderBoard = nil
	for k, v := range a.Arena.State {
		v.URL = k
		arena.PutPlayer(v)
		g.LeaderBoard = append(g.LeaderBoard, v)
	}

	g.Arena = arena
	g.PlayerStateByURL = a.Arena.State
	g.UpdateLeaderBoard()
}

func (g Game) Player(url string) *Player {
	pState := g.PlayerStateByURL[url]
	player := NewPlayerWithUrl(url, pState)
	player.Game = g
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

func (g *Game) UpdateLeaderBoard() {
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
