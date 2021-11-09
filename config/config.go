package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	rootDir    = ".leo"
	dataDir    = "data"
	configFile = "config"
)

type Config struct {
	EthereumRPC string `json:"ethereum_rpc"`
	rootPath    string
}

// Initialize creates a config with default settings if one does not exist.
func Initialize(path string) error {
	rootPath := filepath.Join(path, rootDir)
	confPath := filepath.Join(rootPath, configFile)

	_, err := os.Stat(confPath)
	if err == nil || !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return err
	}

	return NewConfig(path).Save()
}

// NewConfig returns a config with default settings.
func NewConfig(path string) *Config {
	return &Config{
		EthereumRPC: "~/Library/Ethereum/geth.ipc",
		rootPath:    filepath.Join(path, rootDir),
	}
}

// Load loads the config from the root path.
func (c *Config) Load() error {
	path := filepath.Join(c.rootPath, configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

// Save writes the config to the root path.
func (c *Config) Save() error {
	path := filepath.Join(c.rootPath, configFile)

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}

// DataPath returns the path to the data directory.
func (c *Config) DataPath() string {
	return filepath.Join(c.rootPath, dataDir)
}
