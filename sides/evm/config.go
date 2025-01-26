package evm

type Config struct {
	PrivateKey   string
	ContractAddr string `yaml:"contract_address"`
	RPCNode      string `yaml:"rpc_url"`
	ChainID      int64  `yaml:"chain_id"`
}
