package pactus

import (
	"context"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Side struct {
	Client   *Client
	Listener *Listener
	Highway  chan *message.Msg

	Ctx context.Context
}

func NewSide(ctx context.Context, highway chan *message.Msg) (*Side, error) {
	client, err := NewClient(ctx, "") // TODO:read rpc url from config
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client, bypass.PACTUS, highway)

	return &Side{
		Client:   client,
		Listener: listener,
		Ctx:      ctx,
	}, nil
}
