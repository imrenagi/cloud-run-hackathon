package main

import "context"

type State interface {
	Play(ctx context.Context) Move
}

