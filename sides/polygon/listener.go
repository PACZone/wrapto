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

func (l *Listener) Start() error {
	for {
		select {
		case <-l.ctx.Done():
			// state
			return nil
		default:
			if err := l.ProcessOrder(); err != nil {
				return err
			}
		}
	}
}

func (l *Listener) ProcessOrder() error {
	o, err := l.client.Get(*big.NewInt(int64(l.nextOrder)))
	if err != nil {
		return err // TODO: retry 3 time
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
		dbErr := l.db.UpdateOrderStatus(ord.ID, order.FAILED)
		if dbErr != nil {
			return dbErr
		}

		dbErr = l.db.AddLog(&database.Log{
			Actor:       "POLYGON",
			Description: fmt.Sprintf("failed to create order: %s", id),
			Trace:       err.Error(),
		})

		if dbErr != nil {
			return dbErr
		}

		return nil
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

	msg := message.NewMessage(params.MainBypass, l.bypass, ord)
	err = msg.Validate(params.MainBypass)
	if err != nil {
		dbErr := l.db.UpdateOrderStatus(id, order.FAILED)
		if dbErr != nil {
			return dbErr
		}

		dbErr = l.db.AddLog(&database.Log{
			Actor:       "POLYGON",
			Description: "invalid message",
			OrderID:     id,
			Trace:       err.Error(),
		})
		if dbErr != nil {
			return dbErr
		}

		return nil
	}

	l.highway <- msg

	dbErr := l.db.AddLog(&database.Log{
		Actor:       "POLYGON",
		Description: "sent order to highway",
		OrderID:     id,
	})
	if dbErr != nil {
		return dbErr
	}

	return nil
}
