package order

import (
	"fmt"
	"math"

	"github.com/PACZone/wrapto/types/params"
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
	// * unique ID on wrapto system.
	ID string

	// * transaction or contract call that user made on source network.
	TxHash string

	// * address of receiver account on destination network.
	Receiver string

	// * address of sender on source network (account that made bridge transaction).
	Sender string

	// * amount of PAC to be bridged, **including fee**.
	amount float64

	// * status of order on wrapto system.
	Status Status
}

func NewOrder(txHash, sender, receiver string, amount float64) (*Order, error) {
	ID, err := gonanoid.ID(10)
	if err != nil {
		return nil, err // ? panic
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

func (o *Order) Fee() float64 {
	fee := o.amount * params.FeeFraction
	ceiledFee := math.Ceil(fee)

	if ceiledFee <= params.MinimumFee {
		return params.MinimumFee
	}

	if ceiledFee >= params.MaximumFee {
		return params.MaximumFee
	}

	return ceiledFee
}

func (o *Order) Amount() float64 {
	return o.amount - o.Fee()
}

func (o *Order) OriginalAmount() float64 {
	return o.amount
}

func (o *Order) basicCheck() error {
	if o.amount <= params.MinimumFee {
		return BasicCheckError{
			Reason: fmt.Sprintf("amount must be more than %v PAC", params.MinimumFee),
		}
	}

	return nil
}
