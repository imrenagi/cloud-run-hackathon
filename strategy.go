package main

type Strategy interface {
	Play(p *Player) Move
}


