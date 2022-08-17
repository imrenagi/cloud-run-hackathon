package main

type Arena struct {
	Width  int // x start from top left to the right
	Height int // y start from top left to the bottom
}

type Game struct {
	Arena Arena
	// Dimension         []int
	PlayersByPosition map[string]PlayerState
}

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

func NewGame(a ArenaUpdate) Game {
	playersByPosition := make(map[string]PlayerState)
	for _, v := range a.Arena.State {
		playersByPosition[v.GetPosition().String()] = v
	}

	return Game{
		Arena: Arena{
			Width:  a.Arena.Dimensions[0],
			Height: a.Arena.Dimensions[1],
		},
		PlayersByPosition: playersByPosition,
	}
}

func (a ArenaUpdate) GetSelf() PlayerState {
	return a.Arena.State[a.Links.Self.Href]
}
