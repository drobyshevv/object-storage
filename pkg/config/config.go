package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`

	S3 struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
		Region    string `yaml:"region"`
	} `yaml:"s3"`

	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func Load() *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("config.yaml", cfg); err != nil {
		log.Fatalf("error reading config.yaml: %v", err)
	}
	return cfg
}
