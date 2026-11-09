package config

import (
	"errors"
	"os"
)

type Config struct {
	CharacterName string
	AuthToken     string
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{
		CharacterName: os.Getenv("ARTIFACTS_CHARACTER"),
		AuthToken:     os.Getenv("ARTIFACTS_TOKEN"),
	}

	if cfg.CharacterName == "" {
		return Config{}, errors.New("переменная ARTIFACTS_CHARACTER не найдена")
	}
	if cfg.AuthToken == "" {
		return Config{}, errors.New("переменная ARTIFACTS_TOKEN не найдена")
	}

	return cfg, nil
}
