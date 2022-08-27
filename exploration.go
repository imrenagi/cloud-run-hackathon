package main

type Exploration interface {
	Explore(p *Player) Move
}
