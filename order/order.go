package order

import (
	"github.com/matoous/go-nanoid/v2"
)

type Status string
type Type string

const (
	CREATED  Status = "CREATED"
	PENDING  Status = "PENDING"
	COMPLETE Status = "COMPLETE"
	FAILED   Status = "FAILED"
)

const (
	PACTUS_POLYGON Type = "PACTUS_POLYGON"
	POLYGON_PACTUS Type = "POLYGON_PACTUS"
)

type Order struct {
	Id              string
	TxHash          string
	Type            Type
	Receiver        string
	Sender          string
	Amount          uint64
	Status          Status
	Fee             uint32
	Reason          string
	ProcessedHash   string
	DestinationAddr string
	Block           uint32
}

func NewOrder(txHash string, t Type, sender string, amount uint64, destAddr string, block uint32) (*Order, error) {
	id, err := gonanoid.New()
	if err != nil {
		return &Order{}, err
	}

	return &Order{
		Id:              id,
		TxHash:          txHash,
		Type:            t,
		Sender:          sender,
		Amount:          amount,
		DestinationAddr: destAddr,
		Block:           block,
		Status:          CREATED,
		Fee:             calculateFee(),
	}, nil
}

func calculateFee() uint32 {
	return 1_000_000_000
}
