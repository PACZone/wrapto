package config

import (
	"os"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/http"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/sides/evm"
	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Environment string          `yaml:"environment"`
	Logger      logger.Config   `yaml:"logger"`
	Pactus      pactus.Config   `yaml:"pactus"`
	Polygon     evm.Config      `yaml:"polygon"` //! NEW EVM.
	Database    database.Config `yaml:"database"`
	HTTPServer  http.Config     `yaml:"http"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}
	defer file.Close()

	config := &Config{}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	if config.Environment != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, Error{
				reason: err.Error(),
			}
		}
	}

	config.Database.URI = os.Getenv("WRAPTO_MONGO_URI")
	config.Pactus.WalletPass = os.Getenv("WRAPTO_PACTUS_WALLET_PASSWORD")
	config.Polygon.PrivateKey = os.Getenv("WRAPTO_POLYGON_PRIVATE_KEY")

	if err = config.basicCheck(); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	return config, nil
}

func (c *Config) basicCheck() error {
	if c.Environment != "dev" && c.Environment != "prod" {
		return InvalidEnvironmentError{
			Environment: c.Environment,
		}
	}

	return nil
}
