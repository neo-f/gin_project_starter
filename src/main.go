package main

import (
	"context"
	"gin_project_starter/src/controllers"
	"gin_project_starter/src/middlewares"
	"gin_project_starter/src/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	settingLogger()
	watchConfig("configs/application.toml")

	binding.Validator = new(utils.ValidatorV9)
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(middlewares.Logger)
	engine.Use(gin.RecoveryWithWriter(log.Logger))

	controllers.Register(engine)

	server := &http.Server{
		Addr:    viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		Handler: engine,
	}

	log.Info().Msgf("Listening and serving HTTP on %s\n", server.Addr)
	go func() {
		_ = server.ListenAndServe()
	}()
	gracefulShutdown(server)
}

func settingLogger() {
	gin.SetMode(gin.ReleaseMode)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Stack().Caller().Timestamp().Logger()
	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func watchConfig(name string) {
	v := viper.New()
	v.SetConfigFile(name)
	if err := v.ReadInConfig(); err != nil {
		log.Panic().Err(err).Str("name", name).Msg(" reading config file failed")
	}
	mergeInViper(v)
	v.OnConfigChange(func(e fsnotify.Event) {
		mergeInViper(v)
		log.Info().Str("name", e.Name).Msg("config file modified")
	})
	v.WatchConfig()
}

func mergeInViper(v *viper.Viper) {
	if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
		log.Panic().Err(err).Msg("failed reading config files")
	}
}

func gracefulShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed shutdown")
	}
	log.Info().Msg("Server closed")
}
