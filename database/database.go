package database

import (
	"context"
	"fmt"
	"time"

	"github.com/PACZone/wrapto/types/order"
	"go.mongodb.org/mongo-driver/bson"
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

func (db *Database) AddOrder(ord *order.Order) (string, error) {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	_, err := coll.InsertOne(ctx, ord)
	if err != nil {
		return "", err
	}

	return ord.ID, nil
}

func (db *Database) AddLog(orderID, actor, desc, trace string) error {
	coll := db.Client.Database(db.DBName).Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	log := &Log{
		OrderID:     orderID,
		Actor:       actor,
		Description: desc,
		Trace:       trace,
	}

	_, err := coll.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateOrderStatus(id string, status order.Status) error {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{"id": id}

	update := bson.M{"$set": bson.M{"status": status}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no order found with id: %s", id)
	}

	return nil
}

func (db *Database) UpdateOrderDestTxHash(id, hash string) error {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{"id": id}

	update := bson.M{"$set": bson.M{"destination_tx_hash": hash}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no order found with id: %s", id)
	}

	return nil
}

func (db *Database) UpdateOrderReason(id, reason string) error {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{"id": id}

	update := bson.M{"$set": bson.M{"reason": reason}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no order found with id: %s", id)
	}

	return nil
}

func (db *Database) GetOrder(id string) (*order.Order, error) {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{"id": id}

	var result order.Order

	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (db *Database) IsOrderExist(id string) (bool, error) {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{"id": id}

	var result order.Order

	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *Database) UpdatePactusState(height uint32) error {
	coll := db.Client.Database(db.DBName).Collection("state")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"pactus": height,
		},
	}

	result, err := coll.UpdateOne(ctx, bson.M{}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated for height: %d", height)
	}

	return nil
}

func (db *Database) UpdatePolygonState(ordID uint32) error {
	coll := db.Client.Database(db.DBName).Collection("state")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"polygon": ordID,
		},
	}

	result, err := coll.UpdateOne(ctx, bson.M{}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated for ordID: %d", ordID)
	}

	return nil
}

func (db *Database) GetState() (*State, error) {
	coll := db.Client.Database(db.DBName).Collection("state")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	var state State

	err := coll.FindOne(ctx, bson.M{}).Decode(&state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (db *Database) GetLatestOrders(limit int) ([]*order.Order, error) {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	var orders []*order.Order

	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{
		bson.E{"created_at", -1}, //nolint
	}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ord order.Order
		if err := cursor.Decode(&ord); err != nil {
			return nil, err
		}
		orders = append(orders, &ord)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (db *Database) SearchOrders(query string) ([]*order.Order, error) {
	coll := db.Client.Database(db.DBName).Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	var orders []*order.Order

	filter := bson.M{
		"$or": []bson.M{
			{"tx_hash": query},
			{"receiver": query},
			{"sender": query},
			{"dest_network_tx_hash": query},
		},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ord order.Order
		if err := cursor.Decode(&ord); err != nil {
			return nil, err
		}
		orders = append(orders, &ord)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
