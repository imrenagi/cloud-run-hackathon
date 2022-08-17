package main

type State interface {
	Play(g Game) Decision
}

