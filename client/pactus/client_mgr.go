package pactus

import (
	"context"
	"fmt"

	pactus "github.com/PACZone/teleport/client/pactus/gen/go"
)

type Mgr struct {
	ctx     context.Context
	clients []Client
}

func NewClientMgr(ctx context.Context) *Mgr {
	return &Mgr{
		clients: make([]Client, 0),
		ctx:     ctx,
	}
}

func (cm *Mgr) AddClient(c Client) {
	cm.clients = append(cm.clients, c)
}

func (cm *Mgr) GetBlock(height uint32) (*pactus.GetBlockResponse, error) {
	for _, c := range cm.clients {
		txs, err := c.GetBlock(cm.ctx, height)
		if err == nil {
			return txs, nil
		}

		continue
	}

	return nil, fmt.Errorf("can't find block with height %d", height)
}

func (cm *Mgr) GetHeight() (uint32, error) {
	for _, c := range cm.clients {
		lbh, err := c.GetHeight(cm.ctx)
		if err == nil {
			return lbh, nil
		}

		continue
	}

	return 0, fmt.Errorf("can't find last block height")
}
