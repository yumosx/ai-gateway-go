package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	ds "github.com/cohesion-org/deepseek-go"
	"github.com/ecodeclub/ai-gateway-go/ent"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func NewEntClient(lc fx.Lifecycle, cfg Config) (*ent.Client, error) {
	db, err := sql.Open("mysql", cfg.MySQL.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open: %w", err)
	}
	if cfg.MySQL.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	}
	if cfg.MySQL.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	}
	if cfg.MySQL.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql ping: %w", err)
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.MySQL, db)))
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})
	return client, nil
}

func NewRedisClient(lc fx.Lifecycle, cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return rdb.Close()
		},
	})
	return rdb, nil
}

func NewGinEngine() *gin.Engine {
	return gin.Default()
}

func NewDeepSeekClient(cfg Config) *ds.Client {
	return ds.NewClient(cfg.DeepSeek.Token)
}
