package storages

import (
	"context"
	"github.com/spf13/viper"
	"sync"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog/log"
)

var defaultStorage *Storage

type databaseConfig struct {
	Name string
	Dsn  string
}

type Storage struct {
	databases map[string]*pg.DB
	mutex     sync.RWMutex
	connected bool
}

func NewStorage() *Storage {
	s := &Storage{
		databases: map[string]*pg.DB{},
		mutex:     sync.RWMutex{},
		connected: false,
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

func (s *Storage) lazyInit() {
	if s.connected {
		return
	}
	var dbs []databaseConfig
	if err := viper.UnmarshalKey("storage.databases", &dbs); err != nil {
		log.Fatal().Err(err).Msg("database config parse failed")
	}
	for _, db := range dbs {
		if err := s.register(db); err != nil {
			log.Fatal().Err(err).Str("name", db.Name).Msg("database registration failed")
		}
	}
	s.connected = true
}

func Get(name string) *pg.DB { return defaultStorage.Get(name) }
func (s *Storage) Get(name string) *pg.DB {
	s.lazyInit()
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	db, ok := s.databases[name]
	if !ok {
		log.Error().Str("name", name).Msg("database has not registered")
	}
	return db
}

func Close() { defaultStorage.Close() }
func (s *Storage) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for name, conn := range s.databases {
		_ = conn.Close()
		log.Info().Str("name", name).Msg("database closed")
	}
	s.connected = false
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	latency := time.Since(event.StartTime).String()
	query, _ := event.FormattedQuery()
	log.Debug().Str("latency", latency).Msg(query)
	return nil
}

func init() {
	defaultStorage = NewStorage()
}
