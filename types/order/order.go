package order

import (
	"fmt"
	"math"

	"github.com/PACZone/wrapto/types"
	gonanoid "github.com/matoous/go-nanoid"
)

type Status string

const (
	CREATED  Status = "CREATED"
	PENDING  Status = "PENDING"
	COMPLETE Status = "COMPLETE"
	FAILED   Status = "FAILED"
)

type Order struct {
	// * unique ID on wraptor system.
	ID string

	// * transaction or contract call that user made on source network.
	TxHash string

	// * address of receiver account on destination network.
	Receiver string

	// * address of sender on source network (account that made bridge transaction).
	Sender string

	// * amount of PAC to be bridged, **including fee**.
	amount uint64

	// * status of order on wraptor system.
	Status Status

	// * once status got COMPLETE, this will be filled with destination network transaction hash made by wraptor.
	DestNetworkTxHash string

	// * will be filled if order failed.
	Reason string
}

func NewOrder(txHash, sender, receiver string, amount uint64) (*Order, error) {
	ID, err := gonanoid.ID(10)
	if err != nil {
		return nil, err
	}

	ord := &Order{
		ID:       ID,
		TxHash:   txHash,
		Receiver: receiver,
		Sender:   sender,
		amount:   amount,
		Status:   CREATED,
	}

	if err := ord.basicCheck(); err != nil {
		return nil, err
	}

	return ord, nil
}

func (o *Order) Fee() uint64 {
	fee := float64(o.amount) * types.FeeFraction

	ceiledFee := uint64(math.Ceil(fee))

	if ceiledFee <= types.MinimumFee {
		return types.MinimumFee
	}

	if ceiledFee >= types.MaximumFee {
		return types.MaximumFee
	}

	return ceiledFee
}

func (o *Order) Amount() uint64 {
	return o.amount - o.Fee()
}

func (o *Order) basicCheck() error {
	if o.amount <= types.MinimumFee {
		return BasicCheckError{
			Reason: fmt.Sprintf("amount must be more than %d PAC", types.MinimumFee),
		}
	}

	return nil
}
