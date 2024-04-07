package pactus

import (
	"context"
	"fmt"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
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

func (b *Bridge) Start() error {
	logger.Info("starting bridge", "actor", b.bypassName)
	for {
		select {
		case <-b.ctx.Done():
			logger.Info("stopping bridge", "actor", b.bypassName)
			b.wallet.closeWallet()

			return nil
		case msg := <-b.bypass:
			err := b.processMessage(msg)
			if err != nil {
				logger.Error("error while processing message in bridge",
					"actor", b.bypassName, "orderID", msg.Payload.ID)

				return err
			}
		}
	}
}

func (b *Bridge) processMessage(msg message.Message) error {
	logger.Info("new message received for process", "actor", b.bypassName, "orderID", msg.Payload.ID)

	err := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "order received as message", "")
	if err != nil {
		return err
	}

	err = msg.Validate(b.bypassName)
	if err != nil {
		logger.Warn("received an invalid message", "actor", b.bypassName, "err", err)

		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "invalid message", err.Error())
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

	amt, err := amount.NewAmount(payload.Amount() / 1e9) // TODO: FIX ME!!!!!!!!!!!!!!!!!!!
	if err != nil {
		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "failed to cast amount", err.Error())
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(msg.Payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	memo := fmt.Sprintf("bridge from %s to %s by wrapto.app", msg.From, msg.To)

	txID, err := b.wallet.transferTx(payload.Receiver, memo, amt)
	if err != nil {
		logger.Error("can't send transaction to pactus network", "actor", b.bypassName, "err", err, "payload", payload)

		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "bridge failed", err.Error())
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	logger.Info("successful bridge", "actor", b.bypassName, "txID", txID, "orderID", msg.Payload.ID)
	err = b.db.AddLog(msg.Payload.ID, string(b.bypassName), "successfully bridged", txID)
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderStatus(payload.ID, order.COMPLETE)
	if err != nil {
		return err
	}

	return nil
}
