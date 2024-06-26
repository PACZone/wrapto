package database

import (
	"errors"

	"github.com/PACZone/wrapto/types/order"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, DBError{
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
		ID:         ord.ID,
		TxHash:     ord.TxHash,
		Receiver:   ord.Receiver,
		Sender:     ord.Sender,
		Amount:     ord.OriginalAmount(),
		Fee:        ord.Fee(),
		Status:     ord.Status,
		BridgeType: ord.BridgeType,
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

func (db *DB) UpdateOrderDestTxHash(id, hash string) error {
	if err := db.Model(&Order{}).Where("id = ?", id).Update("dest_network_tx_hash", hash).Error; err != nil {
		return DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return nil
}

func (db *DB) UpdateOrderReason(id, reason string) error {
	if err := db.Model(&Order{}).Where("id = ?", id).Update("reason", reason).Error; err != nil {
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

	return db.Save(&state).Error
}

func (db *DB) UpdatePolygonState(polygon uint32) error {
	var state State
	if err := db.First(&state).Error; err != nil {
		return err
	}

	state.Polygon = polygon

	return db.Save(&state).Error
}

func (db *DB) GetState() (*State, error) {
	var state State
	if err := db.First(&state).Error; err != nil {
		return nil, err
	}

	return &state, nil
}

func (db *DB) GetLatestOrders(limit int) ([]*Order, error) {
	var orders []*Order
	if err := db.Order("created_at desc").Limit(limit).Find(&orders).Error; err != nil {
		return nil, DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return orders, nil
}

func (db *DB) SearchOrders(q string) ([]*Order, error) {
	var orders []*Order
	if err := db.Where("tx_hash = ? OR receiver = ? OR sender = ? OR dest_network_tx_hash = ?",
		q, q, q, q).Find(&orders).Error; err != nil {
		return nil, DBError{
			TableName: "Orders",
			Reason:    err.Error(),
		}
	}

	return orders, nil
}
