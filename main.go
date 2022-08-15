package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
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
	// kalau ditembak jangan kabur dulu
	// kalau ditembak dari samping F
	// kalau ditembak dari depan R/L
	// kalau ditembak dari belakang R/L
	// mesti tau siapa yang nembak
	// tapi yang paling penting cari jalan yg kosong
	if player.WasHit {
		commands := []string{"F", "R"}
		rand := rand2.Intn(2)
		return commands[rand]
	}
	//benerin cara nyari lawan soalnya suka muter2 ketika gak ada orang
	// check user di kiri
	playersInFront := player.GetPlayersInDirection(game, player.GetDirection())
	playersInLeft := player.GetPlayersInDirection(game, player.GetDirection().Left())
	playersInRight := player.GetPlayersInDirection(game, player.GetDirection().Right())
	if playersInFront > 0 {
		return "T"
	} else if playersInLeft > 0 {
		return "L"
	} else if playersInRight > 0 {
		return "R"
	} else {
		commands := []string{"F", "R", "L"}
		rand := rand2.Intn(3)
		return commands[rand]
	}
}
