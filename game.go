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

const (
	defaultAttackRange int = 3
)

type GameConfig struct {
	AttackRange int
}

type Game struct {
	Arena             Arena
	PlayersByURL      map[string]PlayerState
	Config            GameConfig
}

func NewGame(a ArenaUpdate) Game {

	width := a.Arena.Dimensions[0]
	height := a.Arena.Dimensions[1]
	arena := NewArena(width, height)

	for _, v := range a.Arena.State {
		arena.PutPlayer(v)
	}

	return Game{
		Arena: arena,
		PlayersByURL:      a.Arena.State,
		Config: GameConfig{
			AttackRange: defaultAttackRange,
		},
	}
}

func (g Game) Player(url string) PlayerState {
	player := g.PlayersByURL[url]
	player.Game = g
	return player
}

func (g Game) GetPlayerByPosition(p Point) (PlayerState, bool) {
	player := g.Arena.Grid[p.Y][p.X].Player
	if player == nil {
		return PlayerState{}, false
	}
	player.Game = g
	return *player, true
}
