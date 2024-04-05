package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Network  string
	Pactus   PactusConfig
	Polygon  PolygonConfig
	Database DatabaseConfig
}

type PactusConfig struct {
	WalletPath string
	WalletPass string
	LockAddr   string
	RPCNode    string
}

type PolygonConfig struct {
	PrivateKey   string
	ContractAddr string
	RPCNode      string
}

type DatabaseConfig struct {
	Path string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Network: os.Getenv("NETWORK"),
		Pactus: PactusConfig{
			WalletPath: os.Getenv("PACTUS_WALLET_PATH"),
			WalletPass: os.Getenv("PACTUS_WALLET_PASSWORD"),
			LockAddr:   os.Getenv("PACTUS_WALLET_ADDRESS"),
			RPCNode:    os.Getenv("PACTUS_RPC"),
		},
		Polygon: PolygonConfig{
			PrivateKey:   os.Getenv("POLYGON_PRIVATE_KEY"),
			ContractAddr: os.Getenv("POLYGON_CONTRACT_ADDRESS"),
			RPCNode:      os.Getenv("POLYGON_RPC"),
		},
		Database: DatabaseConfig{
			os.Getenv("DATABASE_PATH"),
		},
	}

	if err := cfg.basicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) basicCheck() error {
	if !(c.Network == "main") || !(c.Network == "test") {
		return InvalidNetworkError{
			Network: c.Network,
		}
	}

	return nil
}
