package database

import (
	"errors"

	"github.com/PacmanHQ/teleport/order"
	"github.com/glebarez/sqlite"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, errors.New("can't open database")
	}

	if !db.Migrator().HasTable(&Order{}) || !db.Migrator().HasTable(&Listened{}) {
		if err := db.AutoMigrate(
			&Order{},
			&Listened{},
		); err != nil {
			return nil, errors.New("can't auto migrate tables")
		}
	}

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) AddOrder(order order.Order) error {
	return db.Create(&Order{
		Id:              order.Id,
		TxHash:          order.TxHash,
		Type:            order.Type,
		Receiver:        order.Receiver,
		Sender:          order.Sender,
		Amount:          order.Amount,
		Status:          order.Status,
		Fee:             order.Fee,
		Reason:          order.Reason,
		ProcessedHash:   order.ProcessedHash,
		DestinationAddr: order.DestinationAddr,
		Block:           order.Block,
	}).Error
}

func (db *DB) UpdateOrderStatus(id string, status order.Status) error {
	return db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

func (db *DB) UpdateOrderProcessedHashAndReason(id string, processedTxHash string, reason string, s order.Status) error {
	return db.Model(&Order{}).Where("id = ?", id).Updates(map[string]interface{}{"processed_tx_hash": processedTxHash, "reason": reason, "status": s}).Error
}

func (db *DB) GetOrderByTxHash(txHash string) (*Order, error) {
	var order Order
	if err := db.First(&order, "tx_hash = ?", txHash).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (db *DB) CreateListened(network int, last int, txCount int) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	listened := &Listened{
		Id:      id,
		Network: network,
		Last:    last,
		TxCount: txCount,
	}
	return db.Create(listened).Error
}

func (db *DB) GetLastListened() (*Listened, error) {
	var listened Listened
	if err := db.Last(&listened).Error; err != nil {
		return nil, err
	}
	return &listened, nil
}
