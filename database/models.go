package database

import (
	"github.com/PACZone/wrapto/types/order"
	"gorm.io/gorm"
)

type Order struct {
	// * unique ID on wrapto system.
	ID string `gorm:"primaryKey"`

	// * transaction or contract call that user made on source network.
	TxHash string

	// * address of receiver account on destination network.
	Receiver string

	// * address of sender on source network (account that made bridge transaction).
	Sender string

	// * amount of PAC to be bridged, **including fee**.
	Amount float64

	// * fee of order on wrapto system
	Fee float64

	// * status of order on wrapto system.
	Status order.Status

	// * once status got COMPLETE, this will be filled with destination network transaction hash made by wrapto.
	DestNetworkTxHash string

	// * will be filled if order failed.
	Reason string

	gorm.Model

	Logs []Log `gorm:"foreignKey:OrderID"`
}

type Log struct {
	Actor string

	Description string

	Trace string

	gorm.Model

	OrderID string
}
