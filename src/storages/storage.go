package storages

import (
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

type databaseConfig struct {
	Name string
	Dsn  string
}

type Storage struct {
	databases   map[string]*pg.DB
	mutex       sync.RWMutex
	connectOnce sync.Once
}

func NewStorage() *Storage {
	s := &Storage{
		databases:   map[string]*pg.DB{},
		mutex:       sync.RWMutex{},
		connectOnce: sync.Once{},
	}
	return s
}

func (s *Storage) register(conf databaseConfig) error {
	opt, err := pg.ParseURL(conf.Dsn)
	if err != nil {
		return err
	}
	conn := pg.Connect(opt)
	if _, err = conn.Exec("SELECT 1;"); err != nil {
		return err
	}
	log.Info().Str("name", conf.Name).Msg("database connected")
	conn.AddQueryHook(dbLogger{})
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.databases[conf.Name] = conn
	return nil
}

func (s *Storage) Get(name string) *pg.DB {
	s.lazyInit()
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	db, ok := s.databases[name]
	if !ok {
		log.Error().Str("name", name).Msg("database have not registered")
	}
	return db
}

func (s *Storage) lazyInit() {
	s.connectOnce.Do(func() {
		var dbs []databaseConfig
		if err := viper.UnmarshalKey("storage.databases", &dbs); err != nil {
			log.Fatal().Err(err).Msg("database config parse failed")
		}
		for _, db := range dbs {
			if err := s.register(db); err != nil {
				log.Fatal().Err(err).Str("name", db.Name).Msg("database registration failed")
			}
		}
	})
}

func (s *Storage) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for name, conn := range s.databases {
		_ = conn.Close()
		log.Info().Str("name", name).Msg("database closed")
	}
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	q.Data["start"] = time.Now()
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	latency := time.Since(q.Data["start"].(time.Time)).String()
	query, _ := q.FormattedQuery()
	log.Debug().Str("latency", latency).Msg(query)
}
