package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Logger      LoggerConfig
	Pactus      PactusConfig
	Polygon     PolygonConfig
	Database    DatabaseConfig
	HTTPServer  HTTPServerConfig
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

type HTTPServerConfig struct {
	Port string
}

type LoggerConfig struct {
	Filename   string
	LogLevel   string
	Targets    []string
	MaxSize    int
	MaxBackups int
	Compress   bool
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	maxSizeStr := os.Getenv("LOG_MAX_SIZE")
	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		return nil, err
	}

	maxBackupsStr := os.Getenv("LOG_MAX_BACKUPS")
	maxBackups, err := strconv.Atoi(maxBackupsStr)
	if err != nil {
		return nil, err
	}

	compressStr := os.Getenv("LOG_COMPRESS")
	compress, err := strconv.ParseBool(compressStr)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Logger: LoggerConfig{
			Filename:   os.ExpandEnv("LOG_FILENAME"),
			LogLevel:   os.Getenv("LOG_LEVEL"),
			Targets:    strings.Split(os.Getenv("LOG_TARGETS"), ","),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			Compress:   compress,
		},
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

		HTTPServer: HTTPServerConfig{
			Port: os.Getenv("HTTP_PORT"),
		},
	}

	if err := cfg.basicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) basicCheck() error {
	if c.Environment != "dev" && c.Environment != "prod" {
		return InvalidEnvironmentError{
			Environment: c.Environment,
		}
	}

	return nil
}
