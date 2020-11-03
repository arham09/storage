package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db *redis.Client
}

// New creates a new redis storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Create new redis client
	db := redis.NewClient(&redis.Options{
		Network:            cfg.Network,
		Addr:               cfg.Addr,
		Dialer:             cfg.Dialer,
		OnConnect:          cfg.OnConnect,
		Username:           cfg.Username,
		Password:           cfg.Password,
		DB:                 cfg.DB,
		MaxRetries:         cfg.MaxRetries,
		MinRetryBackoff:    cfg.MinRetryBackoff,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		MinIdleConns:       cfg.MinIdleConns,
		MaxConnAge:         cfg.MaxConnAge,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		TLSConfig:          cfg.TLSConfig,
		Limiter:            cfg.Limiter,
	})

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	// Create new store
	return &Storage{
		db: db,
	}
}

// Get value by key
func (s *Storage) Get(key string) ([]byte, error) {
	val, err := s.db.Get(context.Background(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return nil, nil
	}
	return val, nil
}

// Set key with value
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	return s.db.Set(context.Background(), key, val, exp).Err()
}

// Delete key by key
func (s *Storage) Delete(key string) error {
	return s.db.Del(context.Background(), key).Err()
}

// DeleteWild key by prefix
func (s *Storage) DeleteWild(prefix string) error {
	item := s.db.Scan(context.Background(), 0, prefix, 0).Iterator()

	for item.Next(context.Background()) {
		if err := s.db.Del(context.Background(), item.Val()).Err(); err != nil {
			return err
		}
	}

	return item.Err()
}

// Clear all keys
func (s *Storage) Clear() error {
	return s.db.FlushDB(context.Background()).Err()
}
