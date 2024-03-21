package order

import (
	"os"
	"strconv"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type (
	Status string
	Type   string
)

const (
	CREATED  Status = "CREATED"
	PENDING  Status = "PENDING"
	COMPLETE Status = "COMPLETE"
	FAILED   Status = "FAILED"
)

const (
	PACTUS_POLYGON Type = "PACTUS_POLYGON" //nolint
	POLYGON_PACTUS Type = "POLYGON_PACTUS" //nolint
)

type Order struct {
	ID              string
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
	Block           uint32 // todo : update type
}

func NewOrder(txHash string, t Type, sender string, amount uint64,
	destAddr string, block uint32, rec string,
) (*Order, error) {
	id, err := gonanoid.New()
	if err != nil {
		return &Order{}, err
	}

	return &Order{
		ID:              id,
		TxHash:          txHash,
		Type:            t,
		Sender:          sender,
		Amount:          amount,
		DestinationAddr: destAddr,
		Block:           block,
		Status:          CREATED,
		Fee:             getFee(),
		Receiver:        rec,
	}, nil
}

func getFee() uint32 {
	num, err := strconv.ParseUint(os.Getenv("Fee"), 10, 32)
	if err != nil {
		return 0
	}

	return uint32(num)
}

func (o *Order) IsValid() bool {
	return o.Amount > uint64(getFee())
}
