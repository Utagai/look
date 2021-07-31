package config

import (
	"flag"
	"log"
	"os"
)

type BackendType string

const (
	BackendTypeMemory  = "memory"
	BackendTypeMongoDB = "mongodb"
)

type Config struct {
	Source  *os.File
	Backend struct {
		Type    BackendType
		Memory  bool
		MongoDB string
	}
}

func Get() *Config {
	sourcePtr := flag.String("source", "", "the source of data")
	mongodbPtr := flag.String("mongodb", "", "specify the MongoDB connection string URI")

	flag.Parse()

	//// Validate.
	if *sourcePtr == "" {
		log.Fatalf("must specify a source of data")
	}

	//// Set onto Config.
	var cfg Config

	// Source
	source := *sourcePtr
	fi := os.Stdin
	var err error
	if source != "-" {
		fi, err = os.Open(source)
		if err != nil {
			log.Fatalf("failed to open source (%q): %v", source, err)
		}
	}

	cfg.Source = fi

	// Backend type.
	cfg.Backend.Type = BackendTypeMemory
	if *mongodbPtr != "" {
		cfg.Backend.Type = BackendTypeMongoDB
		cfg.Backend.MongoDB = *mongodbPtr
	}

	return &cfg
}
