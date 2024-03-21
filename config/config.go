package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PacLsn PactusListenerConfig
	PolLsn PolygonListenerConfig
	TGBot  TelegramBotConfig
	Wallet WalletConfig
	Fee    uint64
	DBPath string
}

type PactusListenerConfig struct {
	RPCURLS       []string
	BridgeAddress string
}

type PolygonListenerConfig struct {
	ContractAddress string
	RPCURL          string
	PrivateKey      string
}

type WalletConfig struct {
	Path     string
	Address  string
	Password string
}

type TelegramBotConfig struct {
	BotToken string
	ChatID   string
}

func LoadConfig(filePaths ...string) (*Config, error) {
	if err := godotenv.Load(filePaths...); err != nil {
		return nil, err
	}

	f, err := strconv.ParseInt(os.Getenv("Fee"), 10, 64)
	if err != nil {
		return nil, err
	}

	// TODO: should we remove (Pactus | Polygon) prefixes?
	cfg := &Config{
		PacLsn: PactusListenerConfig{
			RPCURLS:       strings.Split(os.Getenv("PACTUS_NODES"), ","),
			BridgeAddress: os.Getenv("PACTUS_BRIDGE_ADDRESS"),
		},
		PolLsn: PolygonListenerConfig{
			ContractAddress: os.Getenv("POLYGON_CONTRACT_ADDRESS"),
			RPCURL:          os.Getenv("POLYGON_RPC"),
			PrivateKey:      os.Getenv("POLYGON_PRIVATE_KEY"),
		},
		TGBot: TelegramBotConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		},
		Wallet: WalletConfig{
			Path:     os.Getenv("WALLET_PATH"),
			Address:  os.Getenv("WALLET_ADDRESS"),
			Password: os.Getenv("WALLET_PASSWORD"),
		},
		Fee:    uint64(f),
		DBPath: os.Getenv("DB_PATH"),
	}

	return cfg, nil
}
