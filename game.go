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

const (
	defaultAttackRange int = 3
)

type GameConfig struct {
	AttackRange int
}

type Game struct {
	Arena            Arena
	PlayerStateByURL map[string]PlayerState
	Players          []PlayerState
	Config           GameConfig
}

func NewGame() Game {
	return Game{
		Config: GameConfig{
			AttackRange: defaultAttackRange,
		},
	}
}

func (g *Game) UpdateArena(a ArenaUpdate) {
	width := a.Arena.Dimensions[0]
	height := a.Arena.Dimensions[1]
	arena := NewArena(width, height)

	g.Players = nil
	for _, v := range a.Arena.State {
		arena.PutPlayer(v)
		g.Players = append(g.Players, v)
	}

	// sort.Sort(byScore(g.Players))
	g.Arena = arena
	g.PlayerStateByURL = a.Arena.State
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

func (g Game) GetPlayerStateByPosition(p Point) (PlayerState, bool) {
	player := g.Arena.Grid[p.Y][p.X].Player
	if player == nil {
		return PlayerState{}, false
	}
	return *player, true
}

func (g Game) GetPlayerByPosition(p Point) *Player {
	pState := g.Arena.Grid[p.Y][p.X].Player
	if pState == nil {
		return nil
	}
	player := NewPlayer(*pState)
	player.Game = g
	return player
}

// ByAge implements sort.Interface based on the Age field.
type byScore []PlayerState

func (a byScore) Len() int           { return len(a) }
func (a byScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
