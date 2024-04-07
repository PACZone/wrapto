package pactus

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Listener struct {
	client          *Client
	db              *database.DB
	nextBlockNumber uint32
	bypassName      bypass.Name
	highway         chan message.Message
	lockAddr        string

	ctx context.Context
}

func newListener(ctx context.Context,
	client *Client, bp bypass.Name, highway chan message.Message,
	startBlock uint32, lockAddr string,
	db *database.DB,
) *Listener {
	return &Listener{
		client:          client,
		db:              db,
		bypassName:      bp,
		highway:         highway,
		nextBlockNumber: startBlock,
		lockAddr:        lockAddr,

		ctx: ctx,
	}
}

func (l *Listener) Start() error {
	logger.Info("starting listener", "actor", l.bypassName)
	for {
		select {
		case <-l.ctx.Done():
			logger.Info("stopping listener", "actor", l.bypassName, "nextBlock", l.nextBlockNumber)
			_ = l.client.Close()

			return nil
		default:
			if err := l.processBlocks(); err != nil {
				logger.Error("can't process block on listener", "actor", l.bypassName, "err", err)

				return err
			}
		}
	}
}

func (l *Listener) processBlocks() error {
	ok, err := l.isEligibleBlock(l.nextBlockNumber)
	if err != nil {
		return err
	}

	if !ok {
		time.Sleep(5 * time.Second)

		return nil
	}

	blk, err := l.client.GetBlock(l.nextBlockNumber)
	if err != nil {
		return err
	}

	l.nextBlockNumber++

	validTxs := l.filterValidTransactions(blk.Txs)

	logger.Info("start processing new block", "actor", l.bypassName, "height", blk.Height)
	for _, tx := range validTxs {
		txHash := hex.EncodeToString(tx.Id)
		sender := tx.GetTransfer().Sender
		amt := float64(tx.GetTransfer().Amount)

		logger.Info("processing new tx", "actor", l.bypassName, "height", blk.Height, "txID", txHash,
			"amount", amt, "sender", sender)

		destInfo, err := ParseMemo(tx.Memo)
		if err != nil {
			logger.Info("invalid memo", "memo", tx.Memo)

			continue
		}

		ord, err := order.NewOrder(txHash, sender, destInfo.Addr, amt)
		if err != nil {
			logger.Error("error while making new order", "actor", l.bypassName, "err", err,
				"height", blk.Height, "txID", txHash)

			dbErr := l.db.UpdateOrderStatus(ord.ID, order.FAILED)
			if dbErr != nil {
				return dbErr
			}
			dbErr = l.db.AddLog("", "PACTUS", fmt.Sprintf("failed to create order: %s", txHash), err.Error())
			if dbErr != nil {
				return dbErr
			}

			continue
		}

		msg := message.NewMessage(destInfo.BypassName, l.bypassName, ord)
		l.highway <- msg

		logger.Info("sending order message to highway", "actor", l.bypassName, "height",
			blk.Height, "txID", txHash, "orderID", ord.ID)
		err = l.db.AddLog(ord.ID, "PACTUS", "sent order to highway", "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Listener) isEligibleBlock(h uint32) (bool, error) {
	lbh, err := l.client.GetLastBlockHeight()
	if err != nil {
		return false, err
	}

	return h < lbh, nil
}

func (l *Listener) filterValidTransactions(txs []*pactus.TransactionInfo) []*pactus.TransactionInfo {
	validTxs := make([]*pactus.TransactionInfo, 0)

	for _, tx := range txs {
		if tx.PayloadType != pactus.PayloadType_TRANSFER_PAYLOAD ||
			tx.GetTransfer().Receiver != l.lockAddr {
			continue
		}

		validTxs = append(validTxs, tx)
	}

	return validTxs
}
