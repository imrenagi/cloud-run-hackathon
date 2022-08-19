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

func NewArena(w, h int) Arena {
	grid := make([][]Point, h)
	for y := range grid {
		grid[y] = make([]Point, w)
		for x := range grid[y] {
			grid[y][x] = Point{x, y}
		}
	}

	return Arena{
		Width:  w,
		Height: h,
		Grid:   grid,
	}
}

type Arena struct {
	Width  int // x start from top left to the right
	Height int // y start from top left to the bottom
	Grid   [][]Point
}

// Traverse with BFS
func (a Arena) Traverse(start Point) []Point {
	var traversedNode []Point

	visited := make([]bool, a.Width*a.Height)
	queue := make([]Point, 0)
	visited[start.Y*a.Width+start.X] = true
	queue = append(queue, start)

	for {
		if len(queue) == 0 {
			break
		}
		pt := queue[0]
		traversedNode = append(traversedNode, pt)
		queue = queue[1:]
		adjancentNodes := a.GetAdjacent(pt)
		for _, n := range adjancentNodes {
			if !visited[n.Y*a.Width+n.X] {
				visited[n.Y*a.Width+n.X] = true
				queue = append(queue, n)
			}
		}
	}

	return traversedNode
}

func (a Arena) GetAdjacent(p Point) []Point {
	var adj []Point
	iterator := [3]int{-1, 0, 1}
	for _, i := range iterator {
		for _, j := range iterator {
			if i != 0 || j != 0 {
				p := Point{X: p.X + i, Y: p.Y + j}
				if p.IsInArena(a) {
					adj = append(adj, p)
				}
			}
		}
	}
	return adj
}

// A Utility Function to check whether given cell (row, col)
// is a valid cell or not.
func (a Arena) IsValid(p Point) bool {
	return p.Y >= 0 && p.Y < a.Height && p.X >= 0 && p.X < a.Width
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
