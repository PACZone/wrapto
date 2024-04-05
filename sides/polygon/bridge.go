package polygon

import (
	"math/big"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge struct {
	client     *Client
	bypassName bypass.Name
	bypass     chan message.Message
}

func NewBridge(bp chan message.Message, bn bypass.Name, c *Client) Bridge {
	return Bridge{
		bypass:     bp,
		bypassName: bn,
		client:     c,
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

func (b *Bridge) ProcessMsg(msg message.Message) error {
	err := msg.BasicCheck(b.bypassName)
	if err != nil {
		return err
	}

	payload := msg.Payload

	amountBigInt := new(big.Int).SetUint64(uint64(payload.Amount()))

	// TODO: log hash
	_, err = b.client.Mint(*amountBigInt, common.HexToAddress(payload.Receiver))
	if err != nil {
		// log
		return err
	}

	return nil
}
