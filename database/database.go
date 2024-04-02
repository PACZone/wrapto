package database

import (
	"errors"

	"github.com/PACZone/teleport/order"
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

func (db *DB) AddOrder(ord *order.Order) error {
	return db.Create(&Order{
		ID:              ord.ID,
		TxHash:          ord.TxHash,
		Type:            ord.Type,
		Receiver:        ord.Receiver,
		Sender:          ord.Sender,
		Amount:          ord.Amount,
		Status:          ord.Status,
		Fee:             ord.Fee,
		Reason:          ord.Reason,
		ProcessedHash:   ord.ProcessedHash,
		DestinationAddr: ord.DestinationAddr,
		Block:           ord.Block,
	}).Error
}

func (db *DB) UpdateOrderStatus(id string, status order.Status) error {
	return db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

func (db *DB) UpdateOrderProcessedHashAndReason(id, processedTxHash, reason string, s order.Status) error {
	return db.Model(&Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"processed_tx_hash": processedTxHash,
		"reason":            reason, "status": s,
	}).Error
}

func (db *DB) GetOrderByTxHash(txHash string) (*Order, error) {
	var ord Order
	if err := db.First(&ord, "tx_hash = ?", txHash).Error; err != nil {
		return nil, err
	}

	return &ord, nil
}

func (db *DB) CreateListened(network, last, txCount int) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	listened := &Listened{
		ID:      id,
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
