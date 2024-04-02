package pactus

import (
	"context"

	pactus "github.com/PACZone/teleport/client/pactus/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn
}

func NewClient(endpoint string) (Client, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Client{}, err
	}

	return Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}, nil
}

func (c *Client) GetBlock(ctx context.Context, h uint32) (*pactus.GetBlockResponse, error) {
	txs, err := c.blockchainClient.GetBlock(ctx, &pactus.GetBlockRequest{
		Height:    h,
		Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS,
	})
	if err != nil {
		return nil, err
	}

	return txs, err
}

func (c *Client) GetHeight(ctx context.Context) (uint32, error) {
	ci, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}

	return ci.LastBlockHeight, nil
}
