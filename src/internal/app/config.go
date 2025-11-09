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
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	Svgicon     string `json:"svgicon"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Badge       string `json:"badge"`

	// Proxy to server so we can cache the urls
	IconProxy string
}

type Config struct {
	Applications []Application `json:"applications"`
}

func (c *Config) Check() bool {
	return true
}

func (c *Config) SetProxyIconLink() {
	for i := range c.Applications {
		if c.Applications[i].Icon != "" {
			c.Applications[i].IconProxy = "/v1/api/iconcache?link=" + c.Applications[i].Icon
		}
	}
}

func (c *Config) GetApplicationsByCategory() map[string][]Application {
	data := make(map[string][]Application)

	for _, app := range c.Applications {
		data[app.Category] = append(data[app.Category], app)
	}
	for category := range data {
		sort.Slice(data[category], func(i, j int) bool {
			return data[category][i].Name < data[category][j].Name
		})
	}

	return data
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

	config.SetProxyIconLink()

	return &config, nil
}
