package database

import (
	"context"
	"time"

	"github.com/PACZone/wrapto/types/order"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DBName       string
	QueryTimeout time.Duration
	Client       *mongo.Client
}

func Connect(cfg Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnectionTimeout)*time.Millisecond)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.URI).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(time.Duration(cfg.ConnectionTimeout) * time.Millisecond).
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilSliceAsEmpty:   true,
		})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.QueryTimeout)*time.Millisecond)
	defer cancel()

	if err := client.Ping(qCtx, nil); err != nil {
		return nil, err
	}

	return &Database{
		Client:       client,
		DBName:       cfg.DBName,
		QueryTimeout: time.Duration(cfg.QueryTimeout) * time.Millisecond,
	}, nil
}

func (db *Database) AddOrder(_ *order.Order) (string, error) {
	return "", nil
}

func (db *Database) AddLog(_, _, _, _ string) error {
	return nil
}

func (db *Database) UpdateOrderStatus(_ string, _ order.Status) error {
	return nil
}

func (db *Database) UpdateOrderDestTxHash(_, _ string) error {
	return nil
}

func (db *Database) UpdateOrderReason(_, _ string) error {
	return nil
}

func (db *Database) GetOrder(_ string) (*order.Order, error) {
	return nil, nil //nolint
}

func (db *Database) IsOrderExist(_ string) (bool, error) {
	return false, nil
}

func (db *Database) UpdatePactusState(_ uint32) error {
	return nil
}

func (db *Database) UpdatePolygonState(_ uint32) error {
	return nil
}

func (db *Database) GetState() (*State, error) {
	return nil, nil //nolint
}

func (db *Database) GetLatestOrders(_ int) ([]*order.Order, error) {
	return nil, nil
}

func (db *Database) SearchOrders(_ string) ([]*order.Order, error) {
	return nil, nil
}
