package main

import (
	"flag"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.yml", "Path to the config in yaml")
}

func main() {
	// Parsing args
	flag.Parse()

	config, err := newConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}
	// Starting API server
	s := newAPIServer(config)

	if err := s.start(); err != nil {
		log.Fatal(err)
	}

}
