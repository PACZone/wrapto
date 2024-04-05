package polygon

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/PACZone/wrapto/types"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/ethereum/go-ethereum/common"
)

type Listener struct {
	client    *Client
	bypass    bypass.Name
	nextOrder uint32
	highway   chan message.Message

	ctx context.Context
}

func newListener(ctx context.Context,
	client *Client, bp bypass.Name, highway chan message.Message, startOrder uint32,
) *Listener {
	return &Listener{
		client:    client,
		bypass:    bp,
		nextOrder: startOrder,
		highway:   highway,
		ctx:       ctx,
	}
}

func (l *Listener) Start() {
	for {
		select {
		case <-l.ctx.Done():
			// state
			return
		default:
			if err := l.ProcessOrder(); err != nil {
				continue
			}
		}
	}
}

func (l *Listener) ProcessOrder() error {
	o, err := l.client.Get(*big.NewInt(int64(l.nextOrder)))
	if err != nil {
		return err
	}

	if o.Sender == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		<-time.After(20 * time.Second)

		return nil
	}

	amt, _ := o.Amount.Float64()
	sender := o.Sender.Hex()
	id := strconv.FormatUint(uint64(l.nextOrder), 10)
	ord, err := order.NewOrder(id, sender, "", amt)
	if err != nil {
		return err
	}
	// fee
	msg := message.NewMessage(types.MainBypass, l.bypass, ord)
	err = msg.BasicCheck(types.MainBypass)
	if err != nil {
		return err
	}

	l.highway <- msg

	return nil
}
