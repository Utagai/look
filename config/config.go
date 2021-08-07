package config

import (
	"flag"
	"fmt"
	"io"
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
	CustomFields *CustomFields
}

// CustomInputParseReader returns an io.Reader given a source reader that
// transforms the given input into JSON based on user-defined custom parse
// rules.
func (cfg *Config) CustomInputParseReader(src io.Reader) (io.Reader, error) {
	return NewCustomFieldsReader(src, cfg.CustomFields)
}

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

	// Custom parsing
	if *customParsePtr {
		parseFields, err := GetCustomFields(flag.Args())
		if err != nil {
			return nil, fmt.Errorf("failed to parse the custom parser regex options: %w", err)
		}

		cfg.CustomFields = parseFields
	} else {
		// Avoid having a nil in the struct if we can, because why not.
		cfg.CustomFields = nil
	}

	return &cfg, nil
}
