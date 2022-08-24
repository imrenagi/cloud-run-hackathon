package main

type Attack struct {
	Player *Player
}

func (a *Attack) Play() Move {
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

		// return a.Player.Walk()

		// TODO add test cases buat ini
		target := a.Player.FindClosestPlayers()
		if target == nil {
			return Throw
		}

		aStar := NewAStar(a.Player.Game.Arena,
			WithIsUnblockFn(func(p Point) bool {
				return !target.CanAttack(p)
			}),
		)
		path, err := aStar.SearchPath(a.Player.GetPosition(), target.GetPosition())
		if err != nil {
			return a.Player.Walk()
		}

		moves := a.Player.RequiredMoves(path, WithOnlyNextMove())
		if len(moves) > 0 {
			return moves[0]
		} else {
			return a.Player.Walk()
		}

		// TODO hindari attack range lawan
		// TODO predict user yg kemungkinan masuk ke attack range kita di depan, kiri atau kanan.. kalau ada stop.
	}
}
