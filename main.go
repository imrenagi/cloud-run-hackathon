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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	ttrace "github.com/imrenagi/cloud-run-hackathon-go/internal/telemetry/trace"
	"github.com/imrenagi/cloud-run-hackathon-go/internal/telemetry/trace/exporter"
)

func init() {
	// zerolog.SetGlobalLevel(zerolog.Disabled)
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

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

	traceProviderCloseFn  []ttrace.CloseFunc
}

const (
	name = "waterfight-server"
)

func NewServer() *Server {
	srv := &Server{
		Router: mux.NewRouter(),
	}
	srv.routes()

	// OTEL_RECEIVER_OTLP_ENDPOINT=localhost:4317
	// TODO do not hardcode
	srv.initGlobalProvider(name, "localhost:4317")

	return srv
}

func (s *Server) routes() {
	s.Router.Use(otelmux.Middleware(name))
	s.Router.Handle("/", s.UpdateArena()).Methods("POST")
	s.Router.Handle("/", s.Healthcheck()).Methods("GET")
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

	for _, closeFn := range s.traceProviderCloseFn {
		go func() {
			err = closeFn(ctxShutDown)
			if err != nil {
				log.Error().Err(err).Msgf("Unable to close trace provider")
			}
		}()
	}

	return err
}

func (s Server) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.player = nil
		fmt.Fprint(w, struct{}{})
	}
}

func (s Server) Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			log.Debug().Msg("Let the battle begin!")
			fmt.Fprint(w, "Let the battle begin!")
			return
		}
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
	// return s.player.Chase(s.player.GetHighestRank())
}

func (s *Server) initGlobalProvider(name, endpoint string) {
	spanExporter := exporter.NewOTLP(endpoint)
	tracerProvider, tracerProviderCloseFn, err := ttrace.NewTraceProviderBuilder(name).
		SetExporter(spanExporter).
		Build()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed initializing the tracer provider")
	}
	s.traceProviderCloseFn = append(s.traceProviderCloseFn, tracerProviderCloseFn)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)
}
