package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := Play(v)
	fmt.Fprint(w, resp)
}

func Play(input ArenaUpdate) (response string) {
	player := input.GetSelf()
	game := NewGame(input)

	// escape // benerin cara escape kalau ada yg nembak. cari jalan yg benar
	// kalau dipinggir jalannya masih salah. gak bisa nyari jalan
	// kalau ada obstacle
	// tapi yang paling penting cari jalan yg kosong
	if player.WasHit {
		return string(player.Escape(game))
	}
	// benerin cara nyari lawan soalnya suka muter2 ketika gak ada orang
	// check user di kiri
	playersInFront := player.GetPlayersInRange(game, player.GetDirection(), 3)
	playersInLeft := player.GetPlayersInRange(game, player.GetDirection().Left(), 3)
	playersInRight := player.GetPlayersInRange(game, player.GetDirection().Right(), 3)
	if len(playersInFront) > 0 {
		return "T"
	} else if len(playersInLeft) > 0 {
		return "L"
	} else if len(playersInRight) > 0 {
		return "R"
	} else {
		return string(player.MoveForward(game))
	}
}
