package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	validator "gopkg.in/validator.v2"
)

type LogTransformer struct {
	KernelLogFile         string `json:"kernel_log_file" validate:"nonzero"`
	ContainerMetadataFile string `json:"container_metadata_file" validate:"nonzero"`
	OutputDirectory       string `json:"output_directory" validate:"nonzero"`
}

func New(path string) (*LogTransformer, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("file does not exist: %s", err)
	}
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %s", err)
	}

	cfg := LogTransformer{}
	err = json.Unmarshal(jsonBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %s", err)
	}

	if err := validator.Validate(cfg); err != nil {
		return &cfg, fmt.Errorf("invalid config: %s", err)
	}

	return &cfg, nil
}
