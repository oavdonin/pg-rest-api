package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct for APIServer
type Config struct {
	Server struct {
		// BindAddre is "ip:port"
		BindAddr string `yaml:"bind_addr"`
		// DatabaseUrl representing PG connection
		DatabaseURL string `yaml:"database_url"`
	} `yaml:"api_server"`
}

// newConfig returns a new decoded Config struct
func newConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
