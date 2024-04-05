package database

import (
	"github.com/PACZone/wrapto/types/order"
	"github.com/glebarez/sqlite"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, DBError{
			DBPath: path,
			Reason: err.Error(),
		}
	}

	if !db.Migrator().HasTable(&Order{}) {
		if err := db.AutoMigrate(
			&Order{},
		); err != nil {
			return nil, DBError{
				DBPath: path,
				Reason: err.Error(),
			}
		}
	}

	if !db.Migrator().HasTable(&Log{}) {
		if err := db.AutoMigrate(
			&Order{},
		); err != nil {
			return nil, DBError{
				DBPath: path,
				Reason: err.Error(),
			}
		}
	}

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) CreateOrder(order *order.Order) (string, error) {
	ord:= Order{
		ID:       order.ID,
		TxHash:   order.TxHash,
		Receiver: order.Receiver,
		Sender:   order.Sender,
		Amount:   order.OriginalAmount(),
		Fee:      order.Fee(),
		Status:   order.Status,
	}
	if err := db.Create(ord).Error; err != nil {
		return "", err
	}
	return order.ID, nil
}

func (db *DB) CreateLog(log *Log) error {
	id, err := gonanoid.ID(10)
	if err != nil {
		return err
	}
	log.ID = id
	return db.Create(log).Error
}

func (db *DB) UpdateOrder(order *Order) error {
	return db.Save(order).Error
}

func (db *DB) GetOrder(id string) (*Order, error) {
	var order Order
	err := db.Preload("Logs").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (db *DB) GetOrderWithLogs(id string) (*Order, error) {
	var order Order
	err := db.Preload("Logs").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (db *DB) GetOrderLogs(orderID string) ([]Log, error) {
	var logs []Log
	err := db.Where("order_id = ?", orderID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
