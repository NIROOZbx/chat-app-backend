package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	Port    string        `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
	Mode    string        `mapstructure:"mode"`
}

type DBConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type RedisConfig struct {
	URL string `mapstructure:"url"`
}

type CloudinaryConfig struct{
	URL string `mapstructure:"url"`
}

type Config struct {
	Server Server      `mapstructure:"server"`
	DB     DBConfig    `mapstructure:"database"`
	Redis  RedisConfig `mapstructure:"redis"`
	Cloudinary CloudinaryConfig `mapstructure:"cloudinary"`
}



func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func LoadConfig(path string) (*Config, error) {

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	rawParams, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}


	expandedParams := os.ExpandEnv(string(rawParams))

	v.AutomaticEnv()
	if err:=v.ReadConfig(strings.NewReader(expandedParams)); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}


	var config Config

	if err :=v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	if err :=config.validate(); err != nil {
       return nil, fmt.Errorf("config validation failed: %w", err)
    }
 
    return &config,nil

}


//validate the config 

func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server.port is required")
	}
	if c.DB.DSN == "" {
		return fmt.Errorf("database.dsn is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("redis.url is required")
	}
	return nil
}
