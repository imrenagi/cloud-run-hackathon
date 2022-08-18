package main

type Escape struct {
	Player PlayerState
}

func (e Escape) Play() Decision {
	front := len(e.Player.FindShooterOnDirection(e.Player.GetDirection()))
	back := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward()))
	left := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Left()))
	right := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Right()))

	// TODO kalau dibelakang ada. jaraknya lebih deket (2 puteran), puter aja.
	// TODO kalau udah kepepet, serang aja.

	if (front > 0 && back == 0 && right == 0 && left == 0) ||
	  (front == 0 && back > 0 && right == 0 && left == 0) {
		// TODO kalau dipinggir, cari jarak terpendek buat kabur. bisa lgsg ke kiri daripada muter ke kanan 3 kali
		// ini mestinya bisa solve kalau pakai path planning
		return TurnRight
	} else {
		return e.Player.Walk()
	}
}