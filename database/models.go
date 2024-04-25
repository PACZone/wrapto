package database

import (
	"time"

	"github.com/PACZone/wrapto/types/order"
	"github.com/pactus-project/pactus/types/amount"
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
	Amount amount.Amount

	// * fee of order on wrapto system
	Fee amount.Amount

	// * status of order on wrapto system.
	Status order.Status

	// * once status got COMPLETE, this will be filled with destination network transaction hash made by wrapto.
	DestNetworkTxHash string

	// * will be filled if order failed.
	Reason string

	BrgType order.BrgType

	CreatedAt time.Time `gorm:"autoCreateTime"`

	Logs []Log `gorm:"foreignKey:OrderID"`

	gorm.Model
}

type Log struct {
	Actor string

	Description string

	Trace string

	OrderID string

	gorm.Model
}

type State struct {
	Pactus uint32

	Polygon uint32

	gorm.Model
}
