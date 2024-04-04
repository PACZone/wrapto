package pactus

import (
	"context"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Side struct {
	client   *Client
	listener *Listener
	bridge   Bridge
	highway  chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context, highway chan message.Message, startBlock uint32, w *Wallet, b chan message.Message) (*Side, error) {
	client, err := NewClient(ctx, "") // TODO:read rpc url from config
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client, bypass.PACTUS, highway, startBlock)

	bridge := NewBridge(w, b, bypass.PACTUS)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,
		bridge:   bridge,

		ctx: ctx,
	}, nil
}
