package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/mode", mode)

	log.Info().Msgf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatal().Msgf("http listen error: %v", err)
}

var answer = Fight

func mode(w http.ResponseWriter, req *http.Request) {
	param1 := req.URL.Query().Get("key")
	answer = Decision(param1)
	log.Debug().Msgf("answer is %s", answer)
	fmt.Fprint(w, answer)
	return
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		log.Debug().Msg("Let the battle begin!")
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

func Play(input ArenaUpdate) Decision {
	game := NewGame(input)
	player := game.Player(input.Links.Self.Href)
	return player.Play()
}
