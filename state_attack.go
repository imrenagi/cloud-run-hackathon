package main

type Attack struct {
	Player PlayerState
}

func (a Attack) Play() Decision {
	playersInFront := a.Player.GetPlayersInRange(a.Player.GetDirection(), 3)
	playersInLeft := a.Player.GetPlayersInRange(a.Player.GetDirection().Left(), 3)
	playersInRight := a.Player.GetPlayersInRange(a.Player.GetDirection().Right(), 3)

	if len(playersInFront) > 0 {
		return Fight
	} else if len(playersInLeft) > 0 {
		return TurnLeft
	} else if len(playersInRight) > 0 {
		return TurnRight
	} else {
		return a.Player.Walk()
	}
}