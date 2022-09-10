package main

import "context"


// TODO exploration can check whether it is worth it to go to that point
type Exploration interface {
	Explore(ctx context.Context, p *Player) Move
}
