package pactus

import (
	"context"
	"fmt"
	"time"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/PACZone/wrapto/types/params"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Listener struct {
	client    *Client
	db        *database.DB
	nextBlock uint32
	bypass    bypass.Name
	highway   chan message.Message
	lockAddr  string

	ctx context.Context
}

func newListener(ctx context.Context,
	client *Client, bp bypass.Name, highway chan message.Message,
	startBlock uint32, lockAddr string,
	db *database.DB,
) *Listener {
	return &Listener{
		client:    client,
		db:        db,
		bypass:    bp,
		highway:   highway,
		nextBlock: startBlock,
		lockAddr:  lockAddr,

		ctx: ctx,
	}
}

func (l *Listener) Start() error {
	for {
		select {
		case <-l.ctx.Done():
			// state
			return nil
		default:
			if err := l.ProcessBlocks(); err != nil {
				return err
			}
		}
	}
}

func (l *Listener) ProcessBlocks() error {
	ok, err := l.isEligibleBlock(l.nextBlock)
	if err != nil {
		return err // TODO: handle errors from client
	}

	if !ok {
		<-time.After(5 * time.Second)

		return nil
	}

	blk, err := l.client.GetBlock(l.nextBlock)
	if err != nil {
		return err // TODO: handle errors from client
	}

	validTxs := l.filterValidTxs(blk.Txs)

	for _, tx := range validTxs {
		dest, err := ParseMemo(tx.Memo)
		if err != nil {
			// log -> db
			continue
		}

		txHash := string(tx.Id)
		sender := tx.GetTransfer().Sender
		amt := float64(tx.GetTransfer().Amount)

		ord, err := order.NewOrder(txHash, sender, dest.Addr, amt)
		if err != nil {
			dbErr := l.db.UpdateOrderStatus(ord.ID, order.FAILED)
			if dbErr != nil {
				return dbErr
			}
			dbErr = l.db.AddLog(&database.Log{
				Actor:       "PACTUS",
				Description: fmt.Sprintf("failed to create order: %s", txHash),
				Trace:       err.Error(),
			})

			if dbErr != nil {
				return dbErr
			}

			continue
		}

		msg := message.NewMessage(dest.BypassName, l.bypass, ord)
		err = msg.Validate(params.MainBypass)
		if err != nil {
			dbErr := l.db.UpdateOrderStatus(ord.ID, order.FAILED)
			if dbErr != nil {
				return dbErr
			}

			dbErr = l.db.AddLog(&database.Log{
				Actor:       "PACTUS",
				Description: "invalid message",
				OrderID:     ord.ID,
				Trace:       err.Error(),
			})
			if dbErr != nil {
				return dbErr
			}

			continue
		}

		l.highway <- msg

		dbErr := l.db.AddLog(&database.Log{
			Actor:       "PACTUS",
			Description: "sent order to highway",
			OrderID:     ord.ID,
		})
		if dbErr != nil {
			return dbErr
		}
	}

	l.nextBlock++

	return nil
}

func (l *Listener) isEligibleBlock(h uint32) (bool, error) {
	lst, err := l.client.GetLastBlockHeight()
	if err != nil {
		return false, err
	}

	return h < lst, nil
}

func (l *Listener) filterValidTxs(txs []*pactus.TransactionInfo) []*pactus.TransactionInfo {
	validTxs := make([]*pactus.TransactionInfo, 0)

	for _, tx := range txs {
		if tx.PayloadType != pactus.PayloadType_TRANSFER_PAYLOAD &&
			tx.GetTransfer().Receiver != l.lockAddr {
			continue
		}

		validTxs = append(validTxs, tx)
	}

	return validTxs
}
