package database

import (
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

func (db *DB) CreateOrder(order *Order) error {
	return db.Create(order).Error
}

func (db *DB) CreateLog(log *Log) error {
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
