package database

import (
	"errors"

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

	if !db.Migrator().HasTable(&Order{}) || !db.Migrator().HasTable(&Log{}) {
		if err := db.AutoMigrate(
			&Order{},
			&Log{},
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
	o := &Order{
		ID:       ord.ID,
		TxHash:   ord.TxHash,
		Receiver: ord.Receiver,
		Sender:   ord.Sender,
		Amount:   ord.OriginalAmount(),
		Fee:      ord.Fee(),
		Status:   ord.Status,
	}
	if err := db.Create(o).Error; err != nil {
		return "", DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return ord.ID, nil
}

func (db *DB) AddLog(orderID, actor, desc, trace string) error {
	log := &Log{
		OrderID:     orderID,
		Actor:       actor,
		Description: desc,
		Trace:       trace,
	}

	if err := db.Create(log).Error; err != nil {
		return DBError{
			TableName: "Logs",
			Reason:    err.Error(),
		}
	}

	return nil
}

func (db *DB) UpdateOrderStatus(id string, status order.Status) error {
	if err := db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return nil
}

func (db *DB) GetOrder(id string) (*Order, error) {
	var ord *Order
	if err := db.First(&ord, "id = ?", id).Error; err != nil {
		return nil, DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return ord, nil
}

func (db *DB) IsOrderExist(txHash string) (bool, error) {
	var ord Order
	err := db.Where("tx_hash = ?", txHash).First(&ord).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return true, nil
}

func (db *DB) GetOrderWithLogs(id string) (*Order, error) {
	var ord *Order
	if err := db.Preload("Logs").First(&ord, "id = ?", id).Error; err != nil {
		return nil, DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return ord, nil
}

func (db *DB) GetOrderLogs(orderID string) ([]Log, error) {
	var logs []Log
	if err := db.Where("order_id = ?", orderID).Find(&logs).Error; err != nil {
		return nil, DBError{
			TableName: "Logs",
			Reason:    err.Error(),
		}
	}

	return logs, nil
}
