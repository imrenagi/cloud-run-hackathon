package main

type Attack struct {
	Player Player
}

func (a Attack) Play() Move {
	playersInFront := a.Player.GetPlayersInRange(a.Player.GetDirection(), 3)
	playersInLeft := a.Player.GetPlayersInRange(a.Player.GetDirection().Left(), 3)
	playersInRight := a.Player.GetPlayersInRange(a.Player.GetDirection().Right(), 3)

	if len(playersInFront) > 0 {
		return Throw
	} else if len(playersInLeft) > 0 {
		return TurnLeft
	} else if len(playersInRight) > 0 {
		return TurnRight
	} else {
		return a.Player.Walk()
	}
}