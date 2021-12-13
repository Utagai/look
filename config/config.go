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
	BackendTypeSubstring = "substring"
	BackendTypeBreeze    = "breeze"
	BackendTypeMongoDB   = "mongodb"
)

// Config represents the configuration for look.
type Config struct {
	Source  *os.File
	Backend struct {
		Type    BackendType
		MongoDB string
	}
	CustomFields *custom.Fields
}

// Get returns the config for the current look process.
func Get() (*Config, error) {
	sourcePtr := flag.String("source", "", "the source of data")
	backendPtr := flag.String("backend", BackendTypeBreeze, "the query backend to use")
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
	cfg.Backend.Type = BackendType(*backendPtr)
	if cfg.Backend.Type == BackendTypeMongoDB && *mongodbPtr != "" {
		cfg.Backend.Type = BackendTypeMongoDB
		cfg.Backend.MongoDB = *mongodbPtr
	}

	// Custom fields.
	// TODO: I think this means the user has to do -custom=true, when we'd ideally
	// have them either omit it for false, or just do -custom for true. Should
	// verify this.
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
