package main

type State interface {
	Play() Decision
}

