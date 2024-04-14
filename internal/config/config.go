package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `yaml:"port"`
	Db   `yaml:"db"`
}

type Db struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

func MustLoad() *Config {
	//cfgName := "C:/Users/xzero/Desktop/other stuff/golang_projects/avito_backend_task_2023-main/internal/config/config.yaml"
	cfgName := "./internal/config/config.yaml"
	if _, err := os.Stat(cfgName); os.IsNotExist(err) {
		log.Fatalf("config file does not exist:%s", cfgName)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(cfgName, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	return &cfg
}
