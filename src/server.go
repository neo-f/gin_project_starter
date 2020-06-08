package main

import (
	"context"
	"errors"
	"gin_project_starter/src/controllers"
	"gin_project_starter/src/middlewares"
	"gin_project_starter/src/storages"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultShutDownTimeout = 60 * time.Second

type Server struct {
	server  *http.Server
	Signals chan os.Signal
}

func NewServer() *Server {
	storages.Connect()
	server := &Server{
		server:  NewAPI(),
		Signals: make(chan os.Signal, 1),
	}
	signal.Notify(server.Signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// register viper hook
	// restart server while config changed
	viper.OnConfigChange(func(e fsnotify.Event) {
		server.Signals <- syscall.SIGHUP
		log.Info().Str("name", e.Name).Msg("config file modified")
	})
	return server
}

func (s *Server) Run() {
	log.Info().Msgf("Listening and serving HTTP on %s\n", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Print(err)
	}
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutDownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Failed shutdown")
	}
	storages.Close()
}

func (s *Server) Restart() {
	s.Shutdown()
	s.server = NewAPI()
	go s.Run()
}

func NewAPI() *http.Server {
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(middlewares.Logger)
	engine.Use(middlewares.ErrorResponder())
	engine.Use(gin.RecoveryWithWriter(log.Logger))
	controllers.Register(engine)

	return &http.Server{Addr: viper.GetString("server.addr"), Handler: engine}
}
