package pactus

import (
	"context"
	"fmt"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/pactus-project/pactus/types/amount"
)

type Bridge struct {
	wallet     *Wallet
	db         *database.DB
	bypassName bypass.Name
	bypass     chan message.Message

	ctx context.Context
}

func newBridge(ctx context.Context, w *Wallet, b chan message.Message, bn bypass.Name, db *database.DB) Bridge {
	return Bridge{
		wallet:     w,
		bypass:     b,
		bypassName: bn,
		db:         db,
		ctx:        ctx,
	}
}

func (b Bridge) Start() error {
	for {
		select {
		case <-b.ctx.Done():
			return nil
		case msg := <-b.bypass:
			err := b.ProcessMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (b Bridge) ProcessMsg(msg message.Message) error {
	err := b.db.AddLog(&database.Log{
		OrderID:     msg.Payload.ID,
		Actor:       "PACTUS",
		Description: "order received as message",
	})
	if err != nil {
		return err
	}

	err = msg.Validate(b.bypassName)
	if err != nil {
		dbErr := b.db.AddLog(&database.Log{
			OrderID:     msg.Payload.ID,
			Actor:       "PACTUS",
			Description: "invalid message",
			Trace:       err.Error(),
		})
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(msg.Payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	payload := msg.Payload

	amt, err := amount.NewAmount(payload.Amount())
	if err != nil {
		dbErr := b.db.AddLog(&database.Log{
			OrderID:     msg.Payload.ID,
			Actor:       "PACTUS",
			Description: "failed to cast amount",
			Trace:       err.Error(),
		})
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(msg.Payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	memo := fmt.Sprintf("bridge from %s to %s by Wrapto.app", msg.From, msg.To)

	_, err = b.wallet.TransferTransaction(payload.Receiver, memo, amt)
	if err != nil {
		dbErr := b.db.AddLog(&database.Log{
			OrderID:     msg.Payload.ID,
			Actor:       "PACTUS",
			Description: "tx failed",
			Trace:       err.Error(),
		})
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	dbErr := b.db.UpdateOrderStatus(payload.ID, order.COMPLETE)
	if dbErr != nil {
		return err
	}

	return nil
}
