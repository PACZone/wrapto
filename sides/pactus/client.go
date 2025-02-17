package pactus

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/types/amount"
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

func NewClient(ctx context.Context, endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint,
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
	var err error
	var blockchainInfo *pactus.GetBlockchainInfoResponse

	for i := 0; i <= 3; i++ {
		blockchainInfo, err = c.blockchainClient.GetBlockchainInfo(c.ctx, &pactus.GetBlockchainInfoRequest{})
		if err == nil {
			return blockchainInfo.LastBlockHeight, nil
		}

		time.Sleep(5 * time.Second)
	}

	return 0, ClientError{
		reason: err.Error(),
	}
}

func (c *Client) GetBlock(h uint32) (*pactus.GetBlockResponse, error) {
	var err error
	var block *pactus.GetBlockResponse

	for i := 0; i <= 3; i++ {
		block, err = c.blockchainClient.GetBlock(c.ctx, &pactus.GetBlockRequest{
			Height:    h,
			Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS,
		})
		if err == nil {
			return block, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, ClientError{
		reason: err.Error(),
	}
}

func (c *Client) GetTransaction(id string) (*pactus.GetTransactionResponse, error) {
	var err error
	var tx *pactus.GetTransactionResponse

	for i := 0; i <= 3; i++ {
		tx, err = c.transactionClient.GetTransaction(c.ctx, &pactus.GetTransactionRequest{
			Id: id,
		})
		if err == nil {
			return tx, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, ClientError{
		reason: err.Error(),
	}
}

// ! PACTUS CHANGE.
func (c *Client) GetTotalLocked() (float64, error) {
	var err error
	var addr *pactus.GetAccountResponse
	var total amount.Amount

	for i := 0; i <= 3; i++ {
		addr, err = c.blockchainClient.GetAccount(c.ctx, &pactus.GetAccountRequest{
			Address: "pc1zgp0x33hehvczq6dggs04gywfqpzl9fea5039gh",
		})
		if err == nil {
			total += amount.Amount(addr.Account.Balance)
			break
		}

		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return 0, ClientError{
			reason: err.Error(),
		}
	}

	for i := 0; i <= 3; i++ {
		addr, err = c.blockchainClient.GetAccount(c.ctx, &pactus.GetAccountRequest{
			Address: "pc1zqyxjatqfhaj3arc727alwl4sa3z8lv2m730eh2",
		})
		if err == nil {
			total += amount.Amount(addr.Account.Balance)
			break
		}

		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return 0, ClientError{
			reason: err.Error(),
		}
	}

	return total.ToPAC(), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
