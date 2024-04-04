package pactus

import (
	"context"

	pactus "github.com/PACZone/wrapto/sides/pactus/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn

	ctx context.Context
}

func NewClient(c context.Context, endpoint string) (*Client, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
		ctx:               c,
	}, nil
}

func (c *Client) GetLastBlockHeight() (uint32, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(c.ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}

	return blockchainInfo.LastBlockHeight, nil
}

func (c *Client) GetBlock(h uint32) (*pactus.GetBlockResponse, error) {
	block, err := c.blockchainClient.GetBlock(c.ctx, &pactus.GetBlockRequest{
		Height:    h,
		Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS,
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
