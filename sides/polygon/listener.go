package polygon

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/PACZone/wrapto/types/params"
	"github.com/ethereum/go-ethereum/common"
)

type Listener struct {
	client    *Client
	db        *database.DB
	bypass    bypass.Name
	nextOrder uint32
	highway   chan message.Message

	ctx context.Context
}

func newListener(ctx context.Context,
	client *Client, bp bypass.Name, highway chan message.Message, startOrder uint32, db *database.DB,
) *Listener {
	return &Listener{
		client:    client,
		db:        db,
		bypass:    bp,
		nextOrder: startOrder,
		highway:   highway,
		ctx:       ctx,
	}
}

func (l *Listener) Start() {
	for {
		select {
		case <-l.ctx.Done():
			// state
			return
		default:
			if err := l.ProcessOrder(); err != nil {
				continue
			}
		}
	}
}

func (l *Listener) ProcessOrder() error {
	o, err := l.client.Get(*big.NewInt(int64(l.nextOrder)))
	if err != nil {
		return err
	}

	if o.Sender == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		<-time.After(20 * time.Second)

		return nil
	}

	amt, _ := o.Amount.Float64()
	sender := o.Sender.Hex()
	id := strconv.FormatUint(uint64(l.nextOrder), 10)
	ord, err := order.NewOrder(id, sender, "", amt)
	if err != nil {
		err = l.db.AddLog(&database.Log{
			Actor:       "POLYGON",
			Description: fmt.Sprintf("failed to create order: %s", id),
			Trace:       err.Error(),
		})

		return err
	}

	id, err = l.db.AddOrder(ord)
	if err != nil {
		return err
	}

	err = l.db.AddLog(&database.Log{
		Actor:       "POLYGON",
		Description: "order created",
		OrderID:     id,
	})
	if err != nil {
		return err
	}

	// fee
	msg := message.NewMessage(params.MainBypass, l.bypass, ord)
	err = msg.Validate(params.MainBypass)
	if err != nil {
		return err
	}

	l.highway <- msg

	err = l.db.AddLog(&database.Log{
		Actor:       "POLYGON",
		Description: "sent order to highway",
		OrderID:     id,
	})
	if err != nil {
		return err
	}

	return nil
}
