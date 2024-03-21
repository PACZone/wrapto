package database

import (
	"github.com/PacmanHQ/teleport/order"
	"gorm.io/gorm"
)

type Order struct {
	ID              string
	TxHash          string
	Type            order.Type
	Receiver        string
	Sender          string
	Amount          uint64
	Status          order.Status
	Fee             uint32
	Reason          string
	ProcessedHash   string
	DestinationAddr string
	Block           uint32

	gorm.Model
}

type Listened struct {
	ID      string
	Network int
	Last    int
	TxCount int

	gorm.Model
}
