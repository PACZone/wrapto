package polygon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge struct {
	client     *Client
	db         *database.DB
	bypassName bypass.Name
	bypass     chan message.Message

	ctx context.Context
}

func newBridge(ctx context.Context, bp chan message.Message,
	bn bypass.Name, client *Client, db *database.DB,
) Bridge {
	return Bridge{
		bypass:     bp,
		db:         db,
		bypassName: bn,
		client:     client,

		ctx: ctx,
	}
}

func (b *Bridge) Start() error {
	logger.Info("starting bridge", "actor", b.bypassName)

	for {
		select {
		case <-b.ctx.Done():
			logger.Info("stopping bridge", "actor", b.bypassName)

			return nil
		case msg := <-b.bypass:
			err := b.processMsg(msg)
			if err != nil {
				logger.Error("error while processing message on bridge", "actor", b.bypassName, "err", err)

				return err
			}
		}
	}
}

func (b *Bridge) processMsg(msg message.Message) error {
	logger.Info("received new message on bridge", "actor", b.bypassName, "orderID", msg.Payload.ID)

	err := b.db.AddLog(msg.Payload.ID, "POLYGON", "order received as message", "")
	if err != nil {
		return err
	}
	err = msg.Validate(b.bypassName)
	if err != nil {
		logger.Warn("received message was invalid", "actor", b.bypassName, "err", err)

		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "invalid message", err.Error())
		if dbErr != nil {
			return err
		}

		dbErr = b.db.UpdateOrderReason(msg.Payload.ID, err.Error())
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
	bigIntAmt := new(big.Int).SetUint64(uint64(payload.Amount()))

	hash, err := b.client.Mint(*bigIntAmt, common.HexToAddress(payload.Receiver))
	if err != nil {
		dbErr := b.db.AddLog(msg.Payload.ID, string(b.bypassName), "tx failed", err.Error())
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

		return err
	}

	logger.Info("wPAC minted successfully", "actor", b.bypassName, "txHahs", hash, "orderID", msg.Payload.ID)

	err = b.db.AddLog(msg.Payload.ID, "POLYGON", fmt.Sprintf("tx success with tx hash: %s", hash), hash)
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderStatus(payload.ID, order.COMPLETE)
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderDestTxHash(payload.ID, hash)
	if err != nil {
		return err
	}

	return nil
}
