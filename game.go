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
	PlayersByPosition map[string]PlayerState
	PlayersByURL      map[string]PlayerState
	Config            GameConfig
}

func NewGame(a ArenaUpdate) Game {

	width := a.Arena.Dimensions[0]
	height := a.Arena.Dimensions[1]
	arena := NewArena(width, height)

	playersByPosition := make(map[string]PlayerState)
	for _, v := range a.Arena.State {
		playersByPosition[v.GetPosition().String()] = v
		arena.PutPlayer(v)
	}

	return Game{
		Arena: arena,
		PlayersByPosition: playersByPosition,
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
	player, ok := g.PlayersByPosition[p.String()]
	if !ok {
		return PlayerState{}, false
	}
	player.Game = g
	return player, true
}
