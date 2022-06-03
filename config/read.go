package config

import (
	"errors"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type BuildFlags struct {
	Ldflags string
	Tags    string
}

type Config struct {
	Global   BuildFlags
	Packages []struct {
		Name       string
		BuildFlags `yaml:",inline"`
	}
}

func Read(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	var newConfig Config
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", configPath, err)
	}

	return &newConfig, nil
}
