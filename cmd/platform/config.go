package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HTTP struct {
		Addr string `toml:"addr"`
	} `toml:"http"`
	MySQL struct {
		DSN             string `toml:"dsn"`
		MaxOpenConns    int    `toml:"max_open_conns"`
		MaxIdleConns    int    `toml:"max_idle_conns"`
		ConnMaxLifetime int    `toml:"conn_max_lifetime"`
	} `toml:"mysql"`
	Redis struct {
		Addr     string `toml:"addr"`
		Password string `toml:"password"`
		DB       int    `toml:"db"`
		PoolSize int    `toml:"pool_size"`
	} `toml:"redis"`
	DeepSeek struct {
		Token string `toml:"token"`
	} `toml:"deepseek"`
}

func NewConfig() (Config, error) {
	path := "config.toml"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		path = v
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}
	return cfg, nil
}
