package pactusclient

import (
	"context"
	"errors"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	blockchainClient  pactus.BlockchainClient
	networkClient     pactus.NetworkClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn
}

type PactusClient struct {
	clients []client
}

func NewPactusClient() *PactusClient {
	return &PactusClient{
		clients: make([]client, 0),
	}
}

func (pc *PactusClient) AddClient(endpoint string) error {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	z := client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		networkClient:     pactus.NewNetworkClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}

	pc.clients = append(pc.clients, z)

	return nil
}

func (pc *PactusClient) GetBlockchainHeight(ctx context.Context) (uint32, error) {
	for _, c := range pc.clients {
		blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
		if err != nil {
			continue
		}

		return blockchainInfo.LastBlockHeight, nil
	}

	return 0, errors.New("unable to get blockchainInfo")
}

func (pc *PactusClient) GetBlock(ctx context.Context, height uint32,
	verbosity pactus.BlockVerbosity,
) (*pactus.GetBlockResponse, error) {
	for _, c := range pc.clients {
		block, err := c.blockchainClient.GetBlock(ctx, &pactus.GetBlockRequest{Height: height, Verbosity: verbosity})
		if err != nil {
			continue
		}

		return block, nil
	}

	return nil, errors.New("unable to get block")
}
