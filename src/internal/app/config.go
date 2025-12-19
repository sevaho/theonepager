package app

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type Application struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"` // Link to icon
	Svgicon     string `json:"svgicon"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Badge       string `json:"badge"`

	// Proxy to server so we can cache the urls
	IconProxy string
}

type Config struct {
	Applications map[string][]Application `json:"applications"`
}

func (c *Config) Check() bool {
	return true
}

func (c *Config) Sort() {
	for category := range c.Applications {
		sort.Slice(c.Applications[category], func(i, j int) bool {
			return c.Applications[category][i].Name < c.Applications[category][j].Name
		})
	}
}

func (c *Config) SetProxyIconLink() {
	for category, applications := range c.Applications {
		for i := range applications {
			if c.Applications[category][i].Icon != "" {
				c.Applications[category][i].IconProxy = "/v1/api/iconcache?link=" + c.Applications[category][i].Icon
			}
		}
	}
}

func parseConfig(configfilepath string) (*Config, error) {
	if configfilepath == "" {
		return nil, errors.New("Empty filepath provided, make sure to provide a config filepath.")
	}

	config := Config{}
	data, err := os.ReadFile(configfilepath)

	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unable to load mail template config: %w", err)
	}

	config.Check()
	config.Sort()
	config.SetProxyIconLink()

	return &config, nil
}
