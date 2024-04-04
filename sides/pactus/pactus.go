package pactus

import (
	"context"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Side struct {
	client   *Client
	listener *Listener
	highway  chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context, highway chan message.Message, startBlock uint32) (*Side, error) {
	client, err := NewClient(ctx, "") // TODO:read rpc url from config
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client, bypass.PACTUS, highway, startBlock)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,

		ctx: ctx,
	}, nil
}
