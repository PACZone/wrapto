package order

import (
	"fmt"
	"time"

	"github.com/PACZone/wrapto/types/params"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/pactus-project/pactus/types/amount"
)

type Status string

const (
	PENDING  Status = "PENDING"
	COMPLETE Status = "COMPLETE"
	FAILED   Status = "FAILED"
)

type BridgeType string

// ! NEW EVM.
const (
	PACTUS_POLYGON BridgeType = "PACTUS_POLYGON" //nolint
	POLYGON_PACTUS BridgeType = "POLYGON_PACTUS" //nolint

	BSC_PACTUS BridgeType = "BSC_PACTUS" //nolint
	PACTUS_BSC BridgeType = "PACTUS_BSC" //nolint

	BASE_PACTUS BridgeType = "BASE_PACTUS" //nolint
	PACTUS_BASE BridgeType = "PACTUS_BASE" //nolint
)

// ! NEW EVM.
func BridgeTypeToMemo(bt BridgeType) string {
	switch bt {
	case POLYGON_PACTUS:
		return "Bridge from Polygon to Pactus by wrapto.app"
	case BSC_PACTUS:
		return "Bridge from Binance Smart Chain to Pactus by wrapto.app"
	case BASE_PACTUS:
		return "Bridge from Base to Pactus by wrapto.app"
	case PACTUS_POLYGON:
		return ""
	case PACTUS_BSC:
		return ""
	case PACTUS_BASE:
		return ""
	}

	return ""
}

type Order struct {
	// * unique ID on Wrapto system.
	ID string `bson:"id"`

	// * transaction or contract call that user made on source network.
	TxHash string `bson:"tx_hash"`

	// * address of receiver account on destination network.
	Receiver string `bson:"receiver"`

	// * address of sender on source network (account that made bridge transaction).
	Sender string `bson:"sender"`

	// * Amount of PAC to be bridged, **including fee**.
	Amount amount.Amount `bson:"amount"`

	// * amount of PAC to be bridged, **including fee**.
	CreatedAt time.Time `bson:"created_at"`

	// * status of order on Wrapto system.
	Status Status `bson:"status"`

	// * once status got COMPLETE, this will be filled with destination network transaction hash made by wrapto.
	DestNetworkTxHash string `bson:"destination_tx_hash"`

	// * will be filled if order failed.
	Reason string `bson:"reason"`

	// * type of bridge.
	BridgeType BridgeType `bson:"bridge_type"`
}

func NewOrder(txHash, sender, receiver string, amt amount.Amount, t BridgeType) (*Order, error) {
	ID, err := gonanoid.ID(10)
	if err != nil {
		return nil, err
	}

	ord := &Order{
		ID:         ID,
		TxHash:     txHash,
		Receiver:   receiver,
		Sender:     sender,
		Amount:     amt,
		Status:     PENDING,
		BridgeType: t,
		CreatedAt:  time.Now(),
	}

	if err := ord.basicCheck(); err != nil {
		return nil, err
	}

	return ord, nil
}

func (o *Order) Fee() amount.Amount {
	fee := o.Amount / params.FeeFraction // 0.5% of amount

	if fee <= params.MinimumFee {
		return params.MinimumFee
	}

	if fee >= params.MaximumFee {
		return params.MaximumFee
	}

	return fee
}

func (o *Order) AmountAfterFee() amount.Amount {
	return o.Amount - o.Fee()
}

func (o *Order) OriginalAmount() amount.Amount {
	return o.Amount
}

func (o *Order) basicCheck() error {
	if o.Amount <= params.MinimumFee {
		return BasicCheckError{
			Reason: fmt.Sprintf("amount must be more than %v PAC", params.MinimumFee),
		}
	}

	return nil
}
