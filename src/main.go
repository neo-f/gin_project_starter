package main

import (
	"gin_project_starter/src/controllers"
	"gin_project_starter/src/services"
	"gin_project_starter/src/utils"
	"os"

	"github.com/gin-gonic/gin/binding"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	settingLogger()
	settingViper()
	binding.Validator = new(utils.ValidatorV9)

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// initialize pprof tools
	if gin.Mode() == gin.DebugMode {
		pprof.Register(engine)
	}
	// register controllers
	controllers.Register(engine)

	// services
	services.Initialize()
	defer services.Close()

	// boot
	if err := engine.Run(viper.GetString("server.host") + ":" + viper.GetString("server.port")); err != nil {
		log.Panic().Err(err).Msg("engine boot failed")
	}
}

func settingViper() {
	// load configs
	viper.SetConfigName("application")
	viper.AddConfigPath("configs/")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Str("err", err.Error()).Msg("Fatal error config file")
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msg("Config file changed:" + e.Name)
	})
	viper.WatchConfig()
}

func settingLogger() {
	// initialize logger
	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Caller().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	}
}
