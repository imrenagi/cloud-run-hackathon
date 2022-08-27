package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}

	api := NewServer()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		oscall := <-ch
		log.Debug().Msgf("system call:%+v", oscall)
		cancel()
	}()

	api.Run(ctx, port)
}

type Server struct {
	Router *mux.Router

	game   Game
	player *Player
}

func NewServer() *Server {
	srv := &Server{
		Router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

func (s *Server) routes() {
	s.Router.Handle("/", s.UpdateArena())
	s.Router.Handle("/reset", s.Reset())
}

// Run ...
func (s *Server) Run(ctx context.Context, port string) error {

	httpS := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.Router,
	}

	log.Info().Msgf("server serving on port %s ", port)

	go func() {
		if err := httpS.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	err := httpS.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatal().Msgf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}

func (s Server) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.player = nil
		fmt.Fprint(w, struct{}{})
	}
}

func (s Server) UpdateArena() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			log.Debug().Msg("Let the battle begin!")
			fmt.Fprint(w, "Let the battle begin!")
			return
		}

		var v ArenaUpdate
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&v); err != nil {
			log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := s.Play(v)
		fmt.Fprint(w, resp)
	}
}

func (s *Server) Play(v ArenaUpdate) Move {
	s.game = NewGame()
	s.game.UpdateArena(v)
	if s.player == nil {
		s.player = s.game.Player(v.Links.Self.Href)
	} else {
		s.game.Update(s.player)
	}
	return s.player.Play()
	// topRank := s.game.LeaderBoard[0]
	// target := s.game.GetPlayerByPosition(Point{topRank.X, topRank.Y})
	// return s.player.Chase(target)
}
