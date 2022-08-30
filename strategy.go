package main

import "context"

type Strategy interface {
	Play(ctx context.Context, p *Player) Move
}


