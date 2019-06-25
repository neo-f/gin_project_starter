package storages

import (
	"sync"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type PgStorage struct {
	once sync.Once
	conn *pg.DB
}

func (s *PgStorage) GetDefault() *pg.DB {
	s.once.Do(func() {
		opt, err := pg.ParseURL(viper.GetString("storage.postgres.default.url"))
		if err != nil {
			log.Panic().Err(err).Msg("parsing postgres url failed")
		}
		s.conn = pg.Connect(opt)
	})
	return s.conn
}
