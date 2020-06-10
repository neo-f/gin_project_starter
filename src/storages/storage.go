package storages

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

var defaultDB string
var defaultStorage *Storage

type DatabaseConfig struct {
	Name string
	Dsn  string
}

type Storage struct {
	config    []DatabaseConfig
	databases map[string]*pg.DB
	mutex     sync.RWMutex
	connected bool
}

func NewStorage(config []DatabaseConfig) *Storage {
	s := &Storage{
		config:    config,
		databases: map[string]*pg.DB{},
		mutex:     sync.RWMutex{},
		connected: false,
	}
	return s
}

func (s *Storage) register(conf DatabaseConfig) error {
	opt, err := pg.ParseURL(conf.Dsn)
	if err != nil {
		return err
	}
	conn := pg.Connect(opt)
	if err := conn.Ping(context.Background()); err != nil {
		return err
	}
	log.Info().Str("name", conf.Name).Msg("Database connected")
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
	for _, db := range s.config {
		if err := s.register(db); err != nil {
			log.Fatal().Err(err).Str("name", db.Name).Msg("Database registration failed")
		}
	}
	s.connected = true
}

func Get(name ...string) *pg.DB {
	return defaultStorage.Get(name...)
}
func (s *Storage) Get(dbName ...string) *pg.DB {
	var name string
	if len(dbName) == 0 {
		name = defaultDB
	} else {
		name = dbName[0]
	}
	s.lazyInit()
	s.mutex.RLock()
	db, ok := s.databases[name]
	if !ok {
		log.Fatal().Str("name", name).Msg("The database is not registered")
	}
	defer s.mutex.RUnlock()
	return db
}

func Close() { defaultStorage.Close() }
func (s *Storage) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for name, conn := range s.databases {
		_ = conn.Close()
		log.Info().Str("name", name).Msg("Database closed")
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
	log.Debug().Str("latency", latency).Msg(string(query))
	return nil
}

func Connect() {
	var dbs []DatabaseConfig
	if err := viper.UnmarshalKey("storage.databases", &dbs); err != nil {
		log.Fatal().Err(err).Msg("Invalid database config")
	}
	if len(dbs) > 0 {
		defaultDB = dbs[0].Name
	}
	defaultStorage = NewStorage(dbs)
}
