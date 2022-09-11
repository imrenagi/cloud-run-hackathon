package main

import "context"

type Exploring struct {
	p *Player
}

func (e Exploring) Play(ctx context.Context) Move {
	return "F"
}

type Attacking struct {

}
