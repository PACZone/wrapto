package pactus

import (
	"context"
	"time"

	pactus "github.com/PACZone/wrapto/sides/pactus/gen/go"
)

type Listener struct {
	Client   *Client
	LstBlock uint32

	Ctx context.Context
}

func NewListener(ctx context.Context, client *Client) *Listener {
	return &Listener{
		Client: client,
		Ctx:    ctx,
	}
}

func (l *Listener) Start() {
	for {
		select {
		case <-l.Ctx.Done():
			//state
			return
		default:
			err := l.ProcessBlocks()
			if err != nil {
				// ?
			}
		}
	}
}

func (l *Listener) ProcessBlocks() error {
	ok, err := l.isEligibleBlock(l.LstBlock)
	if err != nil {
		return err
	}

	if !ok {
		<-time.After(5 * time.Second)
		return nil
	}

	// blk, err := l.Client.GetBlock(l.LstBlock)
	// if err != nil {
	// 	return err
	// }

	// validTxs := l.FilterValidTxs(blk.Txs)

	// for _,tx:=range validTxs{
	// 	// dest,err := parseMemo(tx.Memo)
	// 	// if err!= nil{

	// 	// }


	// }
	// detect destination
	// crete order
	// order basic check
	// create message
	// push to highway

	l.LstBlock++

	return nil

}

func (l *Listener) isEligibleBlock(h uint32) (bool, error) {
	lst, err := l.Client.GetLastBlockHeight()
	if err != nil {
		return false, err
	}

	return h < lst, nil
}

func (l *Listener) FilterValidTxs(txs []*pactus.TransactionInfo) []*pactus.TransactionInfo {
	validTxs := make([]*pactus.TransactionInfo, 0)

	for _, tx := range txs {
		if tx.PayloadType != pactus.PayloadType_TRANSFER_PAYLOAD && tx.GetTransfer().Receiver != "p.bridgeAddr" {
			continue
		}

		validTxs = append(validTxs, tx)
	}

	return validTxs

}
