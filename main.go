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
	"go.opentelemetry.io/otel/sdk/trace"

	ttrace "github.com/imrenagi/cloud-run-hackathon-go/internal/telemetry/trace"
	"github.com/imrenagi/cloud-run-hackathon-go/internal/telemetry/trace/exporter"
)

var tracer = otel.Tracer("github.com/imrenagi/cloud-run-hackathon-go")

func init() {
	if _, logLevel := os.LookupEnv("LOG_LEVEL"); !logLevel {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		switch os.Getenv("LOG_LEVEL") {
		case "DEBUG":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "INFO":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "ERROR":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.Disabled)
		}
	}

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

	traceProviderCloseFn []ttrace.CloseFunc
}

const (
	name = "waterfight-server"
)

func NewServer() *Server {
	srv := &Server{
		Router: mux.NewRouter(),
	}
	srv.routes()
	srv.initTracingProvider(name)

	return srv
}

func (s *Server) routes() {

	if _, ok := os.LookupEnv("TRACING_MODE"); ok {
		s.Router.Use(otelmux.Middleware(name))
	}
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

	log.Info().Msgf("server serving on port %s", port)

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
			log.Debug().Msg("Healthcheck done!")
			fmt.Fprint(w, "Healthcheck done!")
			return
		}
	}
}

func (s Server) UpdateArena() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "Server.UpdateArena")
		defer span.End()
		if r.Method == http.MethodGet {
			log.Debug().Msg("Let the battle begin!")
			fmt.Fprint(w, "Let the battle begin!")
			return
		}

		span.AddEvent("parsing json")
		var v ArenaUpdate
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&v); err != nil {
			log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		span.AddEvent("parsing json completed")
		resp := s.Play(ctx, v)
		fmt.Fprint(w, resp)
	}
}

func (s *Server) Play(ctx context.Context, v ArenaUpdate) Move {
	var opts []GameOption
	if val, ok := os.LookupEnv("PLAYER_MODE"); ok {
		opts = append(opts, WithGameMode(Mode(val)))
	}
	s.game = NewGame(opts...)
	s.game.UpdateArena(ctx, v)
	if s.player == nil {
		s.player = s.game.Player(v.Links.Self.Href)
	} else {
		// TODO kalau URLnya ganti, ignore
		s.game.Update(s.player)
	}

	s.player.UpdateHitCount()

	return s.player.Play(ctx)
}

func (s *Server) initTracingProvider(name string) {
	var spanExporter trace.SpanExporter

	if _, ok := os.LookupEnv("TRACING_MODE"); !ok {
		log.Warn().Msgf("TRACING_MODE is not provided. disabling the tracing")
		return
	}

	switch os.Getenv("TRACING_MODE") {
	case "otlp":
		endpoint, ok := os.LookupEnv("OTLP_RECEIVER_ENDPOINT")
		if !ok {
			log.Fatal().Msgf("OTLP_RECEIVER_ENDPOINT must be provided when enabling otlp trace")
		}
		spanExporter = exporter.NewOTLP(endpoint)
		log.Info().Msgf("using otlp trace exporter")
	case "google_cloud_trace":
		spanExporter = exporter.NewGCP()
		log.Info().Msgf("using google cloud trace exporter")
	case "default":
		log.Warn().Msgf("disabling the tracing")
	}

	if spanExporter == nil {
		return
	}

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
