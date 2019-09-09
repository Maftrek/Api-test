package models

import (
	"crypto/sha256"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	configPath = "config/config.toml"
)

type duration time.Duration

// Config struct
type Config struct {
	ServerOpt  ServerOpt `toml:"ServerOpt"`
	HashSum    []byte
	NatsServer NatsServer `toml:"NatsServer"`
}

func (d *duration) UnmarshalText(text []byte) error {
	temp, err := time.ParseDuration(string(text))
	*d = duration(temp)
	return err
}

// ServerOpt struct
type ServerOpt struct {
	ReadTimeout  duration
	WriteTimeout duration
	IdleTimeout  duration
}

// LoadConfig from path
func LoadConfig(c *Config) {
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return
	}

	c.HashSum = GetHashSum()
}

// GetHashSum of config file
func GetHashSum() []byte {
	paths := []string{
		configPath,
	}
	h := sha256.New()

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil
		}
		defer f.Close()
		if _, err = io.Copy(h, f); err != nil {
			return nil
		}
	}

	return h.Sum(nil)
}
