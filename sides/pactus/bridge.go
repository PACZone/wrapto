package pactus

import (
	"fmt"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/pactus-project/pactus/types/amount"
)

type Bridge struct {
	wallet     *Wallet
	bypassName bypass.Name
	bypass     chan message.Message
}

func newBridge(w *Wallet, b chan message.Message, bn bypass.Name) Bridge {
	return Bridge{
		wallet:     w,
		bypass:     b,
		bypassName: bn,
	}
}

func (b Bridge) Start() {
	for msg := range b.bypass {
		err := b.ProcessMsg(msg)
		if err != nil {
			// TODO: Log
			continue
		}
	}
}

func (b Bridge) ProcessMsg(msg message.Message) error {
	err := msg.BasicCheck(b.bypassName)
	if err != nil {
		return err
	}

	payload := msg.Payload

	amt, err := amount.NewAmount(payload.Amount())
	if err != nil {
		return err
	}

	memo := fmt.Sprintf("bridge from %s to %s by Wraptor.app", msg.From, msg.To)

	_, err = b.wallet.TransferTransaction(payload.Receiver, memo, amt) // TODO: update order
	if err != nil {
		return err
	}

	return nil
}
