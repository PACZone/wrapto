package evm

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pactus-project/pactus/types/amount"
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

func NewPublicClient(rpcURL, cAddr string, chainID int64) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	instance, err := NewWrappedPac(common.HexToAddress(cAddr), client)
	if err != nil {
		return nil, err
	}

	return &Client{
		rpcURL:  rpcURL,
		cAddr:   common.HexToAddress(cAddr),
		chainID: *big.NewInt(chainID),
		wpac:    instance,
	}, nil
}

func newClient(rpcURL, pk, cAddr string, chainID int64) (*Client, error) {
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

func (c *Client) Mint(amt big.Int, to common.Address) (string, error) {
	var err error
	var opts *bind.TransactOpts

	opts, err = bind.NewKeyedTransactorWithChainID(c.pk, &c.chainID)
	if err != nil {
		return "", err
	}
	opts.Value = big.NewInt(0)

	var result *types.Transaction
	for i := 0; i <= 3; i++ {
		result, err = c.wpac.Mint(opts, to, &amt)
		if err == nil {
			return result.Hash().String(), nil
		}

		time.Sleep(5 * time.Second)
	}

	return "", ClientError{
		reason: fmt.Sprintf("can't mint %d wPAC to %s, ::: %v", amt.Int64(), to.String(), err),
	}
}

func (c *Client) Get(orderID big.Int) (BridgeOrder, error) {
	var err error
	var result struct {
		Sender             common.Address
		Amount             *big.Int
		DestinationAddress string
		Fee                *big.Int
	}

	for i := 0; i <= 3; i++ {
		result, err = c.wpac.Bridged(&bind.CallOpts{}, &orderID)
		if err == nil {
			return result, nil
		}

		time.Sleep(5 * time.Second)
	}

	return BridgeOrder{}, ClientError{
		reason: fmt.Sprintf("can't get order %d from contract, ::: %v", orderID.Int64(), err),
	}
}

func (c *Client) TotalSupply() (float64, error) {
	var err error
	var result *big.Int

	for i := 0; i <= 3; i++ {
		result, err = c.wpac.TotalSupply(&bind.CallOpts{})
		if err == nil {
			return amount.Amount(result.Int64()).ToPAC(), nil
		}

		time.Sleep(5 * time.Second)
	}

	return 0, ClientError{
		reason: fmt.Sprintf("can't get total supply %d from contract", err),
	}
}
