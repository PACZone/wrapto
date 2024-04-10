package pactus

import (
	"context"
	"fmt"
	"time"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn

	ctx context.Context
}

func newClient(ctx context.Context, endpoint string) (*Client, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
		ctx:               ctx,
	}, nil
}

func (c *Client) GetLastBlockHeight() (uint32, error) {
	for i := 0; i <= 3; i++ {
		blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(c.ctx, &pactus.GetBlockchainInfoRequest{})
		if err == nil {
			return blockchainInfo.LastBlockHeight, nil
		}

		fmt.Println("AAAAAAAAA")

		time.Sleep(5 * time.Second)
	}

	return 0, ClientError{
		reason: "can't get lastBlockHeight from network",
	}
}

func (c *Client) GetBlock(h uint32) (*pactus.GetBlockResponse, error) {
	for i := 0; i <= 3; i++ {
		block, err := c.blockchainClient.GetBlock(c.ctx, &pactus.GetBlockRequest{
			Height:    h,
			Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS,
		})
		if err == nil {
			return block, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, ClientError{
		reason: "can't get block from network",
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}
