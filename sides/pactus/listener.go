package pactus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/pactus-project/pactus/types/amount"
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
		logger.Info("processing new tx", "actor", l.bypassName, "height", blk.Height, "txID", tx.Id)

		if exist, err := l.checkOrderExist(tx.Id); err != nil {
			return err
		} else if exist {
			logger.Warn("error repetitive transaction", "actor", l.bypassName, "txHash", tx.Id)

			continue
		}

		destInfo, err := ParseMemo(tx.Memo)
		if err != nil {
			logger.Info("invalid memo", "memo", tx.Memo)

			continue
		}

		ord, err := l.createOrder(tx, destInfo.Addr)
		if err != nil {
			if errors.Is(err, database.DBError{}) {
				return err
			}

			continue
		}

		id, err := l.db.AddOrder(ord)
		if err != nil {
			return err
		}

		err = l.db.AddLog(id, string(l.bypassName), "order created", "")
		if err != nil {
			return err
		}

		msg := message.NewMessage(destInfo.BypassName, l.bypassName, ord)

		logger.Info("sending order message to highway", "actor", l.bypassName, "height",
			blk.Height, "txID", tx.Id, "orderID", ord.ID)

		l.highway <- msg

		err = l.db.AddLog(ord.ID, "PACTUS", "sent order to highway", "")
		if err != nil {
			return err
		}
	}

	err = l.db.UpdatePactusState(l.nextBlockNumber - 1)
	if err != nil {
		return err
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

func (l *Listener) checkOrderExist(id string) (bool, error) {
	isExist, err := l.db.IsOrderExist(id)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (l *Listener) createOrder(tx *pactus.TransactionInfo, dest string) (*order.Order, error) { //nolint
	sender := tx.GetTransfer().Sender
	amt := tx.GetTransfer().Amount

	ord, err := order.NewOrder(tx.Id, sender, dest, amount.Amount(amt), order.PACTUS_POLYGON)
	if err != nil {
		logger.Error("error while making new order", "actor", l.bypassName, "err", err, "txID", tx.Id)

		dbErr := l.db.AddLog("", "PACTUS", fmt.Sprintf("failed to create order: %s", tx.Id), err.Error())
		if dbErr != nil {
			return nil, dbErr
		}
	}

	return ord, nil
}
