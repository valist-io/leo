package config

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"

	"github.com/valist-io/leo/p2p"
)

const (
	rootDir    = ".leo"
	dataDir    = "data"
	configFile = "config"
)

type Config struct {
	rootPath string
	// ChainId is the chain id to return in the RPC response.
	ChainId *big.Int `json:"chain_id"`
	// BridgeRPC is the ethereum json rpc to use for bridge data.
	BridgeRPC string `json:"bridge_rpc"`
	// PrivateKey is the base64 encoded libp2p private key.
	PrivateKey string `json:"private_key"`
}

// Init creates a config with default settings if one does not exist.
func Init(path string) (Config, error) {
	cfg := Config{
		ChainId:  big.NewInt(1),
		rootPath: filepath.Join(path, rootDir),
	}

	if err := cfg.Load(); err == nil || !os.IsNotExist(err) {
		return cfg, err
	}
	if err := os.MkdirAll(cfg.rootPath, 0755); err != nil {
		return cfg, err
	}

	priv, _, err := p2p.GenerateKey()
	if err != nil {
		return cfg, err
	}
	cfg.PrivateKey, err = p2p.EncodeKey(priv)
	if err != nil {
		return cfg, err
	}

	return cfg, cfg.Save()
}

// Load loads the config from the root path.
func (c *Config) Load() error {
	data, err := os.ReadFile(c.ConfigPath())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

// Save writes the config to the root path.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(c.ConfigPath(), data, 0666)
}

// ConfigPath returns the path to the config file.
func (c *Config) ConfigPath() string {
	return filepath.Join(c.rootPath, configFile)
}

// DataPath returns the path to the data directory.
func (c *Config) DataPath() string {
	return filepath.Join(c.rootPath, dataDir)
}
