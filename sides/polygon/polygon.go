package polygon

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

func NewSide(ctx context.Context, highway chan message.Message, startOrder uint32, b chan message.Message, rpcURL, pk, cAddr string, chainID int64) (*Side, error) {
	client, err := NewClient(rpcURL, pk, cAddr, chainID)
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client, bypass.POLYGON, highway, startOrder)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,

		ctx: ctx,
	}, nil
}
