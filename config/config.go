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
	Db     *Db               `yaml:"db"`
	Tables map[string]string `yaml:"conditions"`
}

func New() *Config {

	config, err := LookupYml()
	if err != nil {

		log.Fatalf(err.Error())
	}

	GlCfg := &Config{
		Db: &Db{
			Host:     config.Db.Host,
			Port:     config.Db.Port,
			User:     config.Db.User,
			Password: config.Db.Password,
			Name:     config.Db.Name,
			Params: map[string]string{
				"sql_mode": "'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO," +
					"NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'",
			},
			Collation:            "utf8_general_ci",
			AllowNativePasswords: true,
			MaxOpenConnections:   DB_MAX_OPEN_CONNS,
			MaxIdleConnections:   DB_MAX_IDLE_CONNS,
			MaxLifetime:          DB_MAX_LIFETIME,
			MaxIdleLifetime:      DB_MAX_IDLE_LIFETIME,
			MaxPingTimeout:       DB_MAX_PING,
		},
	}

	return GlCfg
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
