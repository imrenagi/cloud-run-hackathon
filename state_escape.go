package main

type Escape struct {
	Player PlayerState
}

func (e Escape) Play() Decision {
	front := len(e.Player.FindShooterOnDirection(e.Player.GetDirection()))
	back := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward()))
	left := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Left()))
	right := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Right()))

	emptyAdjacents := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithEmptyAdjacent())
	if len(emptyAdjacents) == 0 {
		if front > 0 {
			return Fight
		} else if left > 0 {
			return TurnLeft
		} else if right > 0 {
			return TurnRight
		} else {
			return e.Player.Walk()
		}
	}



	// TODO belum bisa nentuin mau belok kanan atau kiri efficiently ketika kabur

	if (front > 0 && back == 0 && right == 0 && left == 0) ||
	  (front == 0 && back > 0 && right == 0 && left == 0) {
		// TODO kalau dipinggir, cari jarak terpendek buat kabur. bisa lgsg ke kiri daripada muter ke kanan 3 kali
		// cari valid adjacent. (kiri atau kanan). Terus cari shortest path ke valid adjacent.
		// ini mestinya bisa solve kalau pakai path planning
		return TurnRight
		// return Fight
	} else {
		return e.Player.Walk()
	}
}
