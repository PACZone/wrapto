package order

import (
	"fmt"

	"github.com/PACZone/wrapto/types/params"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/pactus-project/pactus/types/amount"
)

type Status string

const (
	CREATED  Status = "CREATED"
	PENDING  Status = "PENDING"
	COMPLETE Status = "COMPLETE"
	FAILED   Status = "FAILED"
)

type Order struct {
	// * unique ID on Wrapto system.
	ID string

	// * transaction or contract call that user made on source network.
	TxHash string

	// * address of receiver account on destination network.
	Receiver string

	// * address of sender on source network (account that made bridge transaction).
	Sender string

	// * amount of PAC to be bridged, **including fee**.
	amount amount.Amount

	// * status of order on Wrapto system.
	Status Status
}

func NewOrder(txHash, sender, receiver string, amt amount.Amount) (*Order, error) {
	ID, err := gonanoid.ID(10)
	if err != nil {
		return nil, err // ? panic
	}

	ord := &Order{
		ID:       ID,
		TxHash:   txHash,
		Receiver: receiver,
		Sender:   sender,
		amount:   amt,
		Status:   CREATED,
	}

	if err := ord.basicCheck(); err != nil {
		return nil, err
	}

	return ord, nil
}

func (o *Order) Fee() amount.Amount {
	fee := o.amount / params.FeeFraction // 0.5% of amount

	if fee <= params.MinimumFee {
		return params.MinimumFee
	}

	if fee >= params.MaximumFee {
		return params.MaximumFee
	}

	return fee
}

func (o *Order) Amount() amount.Amount {
	return o.amount - o.Fee()
}

func (o *Order) OriginalAmount() amount.Amount {
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
