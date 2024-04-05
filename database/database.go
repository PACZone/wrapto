package database

import (
	"github.com/PACZone/wrapto/types/order"
	"github.com/glebarez/sqlite"
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

func (db *DB) AddOrder(ord *order.Order) (string, error) {
	o := Order{
		ID:       ord.ID,
		TxHash:   ord.TxHash,
		Receiver: ord.Receiver,
		Sender:   ord.Sender,
		Amount:   ord.OriginalAmount(),
		Fee:      ord.Fee(),
		Status:   ord.Status,
	}
	if err := db.Create(o).Error; err != nil {
		return "", err
	}

	return ord.ID, nil
}

func (db *DB) AddLog(log *Log) error {
	return db.Create(log).Error
}

func (db *DB) UpdateOrder(ord *Order) error {
	return db.Save(ord).Error
}

func (db *DB) GetOrder(id string) (*Order, error) {
	var ord Order
	err := db.Preload("Logs").First(&ord, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &ord, nil
}

func (db *DB) GetOrderWithLogs(id string) (*Order, error) {
	var ord Order
	err := db.Preload("Logs").First(&ord, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &ord, nil
}

func (db *DB) GetOrderLogs(orderID string) ([]Log, error) {
	var logs []Log
	err := db.Where("order_id = ?", orderID).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}
