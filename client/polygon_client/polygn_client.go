package polygonClient

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PolygonClient struct {
	rpcUrl  string
	pk      *ecdsa.PrivateKey
	cAddr   common.Address
	chainId big.Int
	wpac    *WrappedPac
}

type bridgeOrder struct {
	Sender             common.Address
	Amount             *big.Int
	DestinationAddress string
	Fee                *big.Int
}

func NewPolygonClient(rpcUrl string, pk string, cAddr string, chainId int64) (*PolygonClient, error) {

	client, err := ethclient.Dial(rpcUrl)
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

	return &PolygonClient{
		rpcUrl:  rpcUrl,
		pk:      privateKey,
		cAddr:   common.HexToAddress(cAddr),
		chainId: *big.NewInt(chainId),
		wpac:    instance,
	}, nil
}

func (p *PolygonClient) Mint(amount big.Int, to common.Address) (string, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(p.pk, &p.chainId)
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

func (p *PolygonClient) GetOrder(orderId big.Int) (bridgeOrder, error) {
	result, err := p.wpac.Bridged(&bind.CallOpts{}, &orderId)
	if err != nil {
		return bridgeOrder{}, nil
	}

	return result, nil
}
