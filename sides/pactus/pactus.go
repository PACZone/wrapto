package pactus

import "context"

type Side struct {
	Client   *Client
	Listener *Listener

	Ctx context.Context
}

func NewPactusSide(ctx context.Context) (*Side, error) {

	client, err := NewClient(ctx, "")
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client)

	return &Side{
		Client:   client,
		Listener: listener,
		Ctx:      ctx,
	}, nil
}
