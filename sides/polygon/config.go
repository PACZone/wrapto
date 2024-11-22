package polygon

type Config struct {
	PrivateKey   string
	ContractAddr string `yaml:"contract_address"`
	RPCNode      string `yaml:"rpc_url"`
}
