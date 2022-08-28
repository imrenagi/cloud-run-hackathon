package main

import "context"

type Exploration interface {
	Explore(ctx context.Context, p *Player) Move
}
