package pactus

import (
	"context"
	"fmt"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
)

type Bridge struct {
	client     *Client
	wallet     *Wallet
	db         *database.Database
	bypassName bypass.Name
	bypass     chan message.Message

	ctx context.Context
}

func newBridge(ctx context.Context, w *Wallet, b chan message.Message,
	bn bypass.Name, db *database.Database, client *Client,
) Bridge {
	return Bridge{
		client:     client,
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

	err := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "order received", fmt.Sprintf("sender: %s", msg.From))
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
	memo := order.BridgeTypeToMemo(msg.Payload.BridgeType)

	txID, tx, err := b.wallet.transferTx(payload.Receiver, memo, payload.OriginalAmount())
	if err != nil {
		logger.Error("can't send transaction to Pactus network", "actor", b.bypassName, "err", err, "payload", payload)

		trace := fmt.Sprintf("error: %s\n", err.Error())
		if tx != nil {
			trace += fmt.Sprintf("\ntx id: %s, tx hex: %s, tx fee: %s, tx amount: %s\n",
				tx.ID().String(), tx.String(), tx.Fee().String(), tx.Payload().Value().String())
		}

		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "bridge failed", trace)
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderReason(msg.Payload.ID, err.Error())
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderStatus(payload.ID, order.FAILED)
		if dbErr != nil {
			return err
		}

		return nil
	}

	_, err = b.client.GetTransaction(txID)
	if err != nil {
		logger.Error("sending tx to pactus network went wrong", "actor", b.bypassName, "err", err, "payload", payload)

		trace := fmt.Sprintf("error: %s\n", err.Error())
		if tx != nil {
			trace += fmt.Sprintf("\ntx id: %s, tx hex: %s, tx fee: %s, tx amount: %s\n",
				tx.ID().String(), tx.String(), tx.Fee().String(), tx.Payload().Value().String())
		}

		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "bridge failed", trace)
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderReason(msg.Payload.ID, err.Error())
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

	trace := txID
	if tx != nil {
		trace += fmt.Sprintf("\ntx id: %s, tx hex: %s, tx fee: %s, tx amount: %s\n",
			tx.ID().String(), tx.String(), tx.Fee().String(), tx.Payload().Value().String())
	}

	err = b.db.AddLog(msg.Payload.ID, string(b.bypassName), "successfully bridged", trace)
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderStatus(payload.ID, order.COMPLETE)
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderDestTxHash(payload.ID, txID)
	if err != nil {
		return err
	}

	return nil
}
