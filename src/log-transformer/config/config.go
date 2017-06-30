package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LogTransformer struct {
	InputFile       string `json:"input_file"`
	OutputDirectory string `json:"output_directory"`
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

	return &cfg, nil
}
