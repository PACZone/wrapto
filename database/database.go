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

	if !db.Migrator().HasTable(&Order{}) || !db.Migrator().HasTable(&Log{}) || !db.Migrator().HasTable(&State{}) {
		if err := db.AutoMigrate(
			&Order{},
			&Log{},
			&State{},
		); err != nil {
			return nil, DBError{
				DBPath: path,
				Reason: err.Error(),
			}
		}

		defaultState := &State{
			Pactus:  638364,
			Polygon: 1,
		}
		if err := db.Create(defaultState).Error; err != nil {
			return nil, DBError{
				TableName: "State",
				Reason:    err.Error(),
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

func (db *DB) UpdatePactusState(pactus uint32) error {
	var state State
	if err := db.First(&state).Error; err != nil {
		return err
	}

	state.Pactus = pactus
	if err := db.Save(&state).Error; err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdatePolygonState(polygon uint32) error {
	var state State
	if err := db.First(&state).Error; err != nil {
		return err
	}

	state.Polygon = polygon
	if err := db.Save(&state).Error; err != nil {
		return err
	}

	return nil
}

func (db *DB) GetState() (*State, error) {
	var state State
	if err := db.First(&state).Error; err != nil {
		return nil, err
	}

	return &state, nil
}
