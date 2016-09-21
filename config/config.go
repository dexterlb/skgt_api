package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database Database `toml:"database"`
	Server   Server   `toml:"server"`
	Parser   Parser   `toml:"parser"`
}

type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"db_name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	SSLMode  string `toml:"ssl_mode"`
}

type Server struct {
	ListenAddress string `toml:"listen_address"`
}

type Parser struct {
	ParallelRequests int `toml:"parallel_requests"`
}

func (db *Database) URN() string {
	var parameters []string

	addURNParameter(&parameters, "host='%s'", db.Host)
	addURNParameter(&parameters, "port='%d'", db.Port)
	addURNParameter(&parameters, "dbname='%s'", db.Name)
	addURNParameter(&parameters, "user='%s'", db.User)
	addURNParameter(&parameters, "password='%s'", db.Password)
	addURNParameter(&parameters, "sslmode='%s'", db.SSLMode)

	return strings.Join(parameters, " ")
}

func addURNParameter(parameters *[]string, format string, parameter interface{}) {
	if parameter != reflect.Zero(reflect.TypeOf(parameter)).Interface() {
		*parameters = append(*parameters, fmt.Sprintf(format, parameter))
	}
}

// Load loads a configuration file
func Load(filename string) (*Config, error) {
	config := &Config{}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	md, err := toml.DecodeReader(f, &config)
	if err != nil {
		return nil, err
	}
	undecoded := md.Undecoded()

	if len(undecoded) > 0 {
		return nil, fmt.Errorf("unknown config values: %v", undecoded)
	}

	return config, nil
}
