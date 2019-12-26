package main

import (
	"gin_project_starter/src/utils"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
)

const configFile = "application.toml"

func main() {
	binding.Validator = new(utils.ValidatorV10)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Stack().Caller().Timestamp().Logger()
	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg(" reading config file failed")
	}
	viper.WatchConfig()

	server := NewServer()
	go server.Run()
	for sig := range server.Signals {
		switch sig {
		case syscall.SIGHUP:
			log.Warn().Msg("Restarting Server ...")
			server.Restart()
		default:
			log.Warn().Msg("Shutting down Server ...")
			server.Shutdown()
			log.Warn().Msg("Server closed")
			return
		}
	}
}
