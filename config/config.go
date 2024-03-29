package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/utagai/look/config/custom"
)

// BackendType is the type of backend to use. The type of backend represents not
// only where the queried data resides, but also the method through which it can
// be queried.
type BackendType string

// The various backend types.
const (
	BackendTypeMemory  = "memory"
	BackendTypeMongoDB = "mongodb"
)

// Config represents the configuration for look.
type Config struct {
	Source  *os.File
	Backend struct {
		Type BackendType
		Memory  bool
		MongoDB string
	}
	CustomFields *custom.Fields
}

// Get returns the config for the current look process.
func Get() (*Config, error) {
	sourcePtr := flag.String("source", "", "the source of data")
	mongodbPtr := flag.String("mongodb", "", "specify the MongoDB connection string URI")
	customParsePtr := flag.Bool("custom-parse", false, "enables custom parsing of the input into JSON")

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

	// Custom fields.
	if *customParsePtr {
		parseFields, err := custom.ParseFields(flag.Args())
		if err != nil {
			return nil, fmt.Errorf("failed to parse the custom parser regex options: %w", err)
		}

		cfg.CustomFields = parseFields
	} else {
		cfg.CustomFields = nil
	}

	return &cfg, nil
}
