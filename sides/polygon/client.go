package polygon

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	rpcURL  string
	pk      *ecdsa.PrivateKey
	cAddr   common.Address
	chainID big.Int
	wpac    *WrappedPac
}

type BridgeOrder struct {
	Sender             common.Address
	Amount             *big.Int
	DestinationAddress string
	Fee                *big.Int
}

func NewClient(rpcURL, pk, cAddr string, chainID int64) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	instance, err := NewWrappedPac(common.HexToAddress(cAddr), client)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	return &Client{
		rpcURL:  rpcURL,
		pk:      privateKey,
		cAddr:   common.HexToAddress(cAddr),
		chainID: *big.NewInt(chainID),
		wpac:    instance,
	}, nil
}

func (p *Client) Mint(amount big.Int, to common.Address) (string, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(p.pk, &p.chainID)
	if err != nil {
		return "", err
	}
	opts.Value = big.NewInt(0)

	result, err := p.wpac.Mint(opts, to, &amount)
	if err != nil {
		return "", err
	}

	return result.Hash().String(), nil
}

func (p *Client) Get(orderID big.Int) (BridgeOrder, error) {
	result, err := p.wpac.Bridged(&bind.CallOpts{}, &orderID)
	if err != nil {
		return BridgeOrder{}, err
	}

	return result, nil
}
