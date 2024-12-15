package pactus

type Config struct {
	WalletPath string `yaml:"wallet_path"`
	WalletPass string
	LockAddr   string `yaml:"lock_address"`
	WalletAddr string `yaml:"wallet_address"`
	RPCNode    string `yaml:"rpc_url"`
}
