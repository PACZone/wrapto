package pactusClient

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PactusClient struct {
	blockchainClient  pactus.BlockchainClient
	networkClient     pactus.NetworkClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn
}

func NewPactusClient(endpoint string) (*PactusClient, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &PactusClient{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		networkClient:     pactus.NewNetworkClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}, nil
}

func (c *PactusClient) GetBlockchainHeight(ctx context.Context) (uint32, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}
	return blockchainInfo.LastBlockHeight, nil
}

func (c *PactusClient) GetBlock(ctx context.Context, height uint32, verbosity pactus.BlockVerbosity) (*pactus.GetBlockResponse, error) {
	block, err := c.blockchainClient.GetBlock(ctx, &pactus.GetBlockRequest{Height: height, Verbosity: verbosity})
	if err != nil {
		return &pactus.GetBlockResponse{}, err
	}
	return block, nil
}
