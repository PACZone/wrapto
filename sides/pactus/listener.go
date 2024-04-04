package pactus

import (
	"context"
	"time"

	pactus "github.com/PACZone/wrapto/sides/pactus/gen/go"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
)

type Listener struct {
	client             *Client
	LastProcessedBlock uint32
	bypass             bypass.Name
	highway            chan *message.Msg

	ctx context.Context
}

func NewListener(ctx context.Context, client *Client, bypass bypass.Name, highway chan *message.Msg) *Listener {
	return &Listener{
		client:  client,
		bypass:  bypass,
		highway: highway,
		ctx:     ctx,
	}
}

func (l *Listener) Start() {
	for {
		select {
		case <-l.ctx.Done():
			// state
			return
		default:
			if err := l.ProcessBlocks(); err != nil {
				continue
			}
		}
	}
}

func (l *Listener) ProcessBlocks() error {
	ok, err := l.isEligibleBlock(l.LastProcessedBlock)
	if err != nil {
		return err // TODO: handle errors from client
	}

	if !ok {
		<-time.After(5 * time.Second)
		return nil
	}

	blk, err := l.client.GetBlock(l.LastProcessedBlock)
	if err != nil {
		return err // TODO: handle errors from client
	}

	validTxs := l.FilterValidTxs(blk.Txs)

	for _, tx := range validTxs {
		dest, err := parseMemo(tx.Memo)
		if err != nil {
			// log -> db
			continue
		}

		txHash := string(tx.Id)
		sender := tx.GetTransfer().Sender
		amt := uint64(tx.GetTransfer().Amount)

		ord, err := order.NewOrder(txHash, sender, dest.addr, amt)
		if err != nil {
			// log -> db
			continue
		}

		msg := message.NewMsg(dest.name, l.bypass, ord)

		l.highway <- msg

	}

	l.LastProcessedBlock++

	return nil
}

func (l *Listener) isEligibleBlock(h uint32) (bool, error) {
	lst, err := l.client.GetLastBlockHeight()
	if err != nil {
		return false, err
	}

	return h < lst, nil
}

func (l *Listener) FilterValidTxs(txs []*pactus.TransactionInfo) []*pactus.TransactionInfo {
	validTxs := make([]*pactus.TransactionInfo, 0)

	for _, tx := range txs {
		if tx.PayloadType != pactus.PayloadType_TRANSFER_PAYLOAD && tx.GetTransfer().Receiver != "LOCKED_ADDRESS" { // TODO:read LOCKED ADDRESS from config
			continue
		}

		validTxs = append(validTxs, tx)
	}

	return validTxs
}
