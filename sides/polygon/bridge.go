package polygon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/PACZone/wrapto/database"
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

func newBridge(ctx context.Context, bp chan message.Message, bn bypass.Name, c *Client, db *database.DB) Bridge {
	return Bridge{
		bypass:     bp,
		db:         db,
		bypassName: bn,
		client:     c,

		ctx: ctx,
	}
}

func (b Bridge) Start() error {
	for {
		select {
		case <-b.ctx.Done():
			// state
			return nil
		case msg := <-b.bypass:
			err := b.ProcessMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (b *Bridge) ProcessMsg(msg message.Message) error {
	err := b.db.AddLog(&database.Log{
		OrderID:     msg.Payload.ID,
		Actor:       "POLYGON",
		Description: "order received as message",
	})
	if err != nil {
		return err
	}
	err = msg.Validate(b.bypassName)
	if err != nil {
		dbErr := b.db.AddLog(&database.Log{
			OrderID:     msg.Payload.ID,
			Actor:       "POLYGON",
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

	amountBigInt := new(big.Int).SetUint64(uint64(payload.Amount()))

	hash, err := b.client.Mint(*amountBigInt, common.HexToAddress(payload.Receiver))
	if err != nil {
		dbErr := b.db.AddLog(&database.Log{
			OrderID:     msg.Payload.ID,
			Actor:       "POLYGON",
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

		return err
	}

	err = b.db.AddLog(&database.Log{
		OrderID:     msg.Payload.ID,
		Actor:       "POLYGON",
		Description: fmt.Sprintf("tx success with tx hash: %s", hash),
	})
	if err != nil {
		return err
	}

	err = b.db.UpdateOrderStatus(payload.ID, order.COMPLETE)
	if err != nil {
		return err
	}

	return nil
}
