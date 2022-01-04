package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	DEFAUL_DB_PORT       = "3306"
	DB_MAX_OPEN_CONNS    = 10
	DB_MAX_IDLE_CONNS    = 10
	DB_MAX_LIFETIME      = 420
	DB_MAX_IDLE_LIFETIME = 420
	DB_MAX_PING          = 5
)

type Db struct {
	Host                 string `yaml:"DB_HOST"`
	Port                 string `yaml:"DB_PORT"`
	User                 string `yaml:"DB_USER"`
	Password             string `yaml:"DB_PASS"`
	Name                 string `yaml:"DB_NAME"`
	Params               map[string]string
	Collation            string
	AllowNativePasswords bool
	MaxOpenConnections   int
	MaxIdleConnections   int
	MaxLifetime          int
	MaxIdleLifetime      int
	MaxPingTimeout       int
}

type Config struct {
	Db            *Db               `yaml:"db"`
	Tables        map[string]string `yaml:"tables"`
	InternalCount int               `yaml:"internal_count"`
	ExternalCount int               `yaml:"external_count"`
	Command       string            `yaml:"command"`
}

func New() *Config {

	config, err := LookupYml()
	if err != nil {

		log.Fatalf(err.Error())
	}

	config.Db.Params = map[string]string{
		"sql_mode": "'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO," +
			"NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'",
	}
	config.Db.Collation = "utf8_general_ci"
	config.Db.AllowNativePasswords = true
	config.Db.MaxOpenConnections = DB_MAX_OPEN_CONNS
	config.Db.MaxIdleConnections = DB_MAX_IDLE_CONNS
	config.Db.MaxLifetime = DB_MAX_LIFETIME
	config.Db.MaxIdleLifetime = DB_MAX_IDLE_LIFETIME
	config.Db.MaxPingTimeout = DB_MAX_PING

	return config
}

func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func LookupYml() (*Config, error) {
	var config Config

	configData, err := ioutil.ReadFile("config.yml")
	if err != nil {

		return nil, fmt.Errorf("yml file not found, %w", err)
	}

	if err := config.Parse(configData); err != nil {

		return nil, fmt.Errorf("can't parse yml file, %w", err)
	}

	return &config, nil
}
