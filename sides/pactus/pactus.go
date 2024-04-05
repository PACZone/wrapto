package pactus

import (
	"context"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/pactus-project/pactus/crypto"
)

type Side struct {
	client   *Client
	listener *Listener
	bridge   Bridge
	highway  chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context,
	highway chan message.Message, startBlock uint32,
	b chan message.Message, net string, cfg config.PactusConfig,
) (*Side, error) {
	if net == "test" {
		crypto.AddressHRP = "tpc"
	}

	client, err := NewClient(ctx, cfg.RPCNode)
	if err != nil {
		return nil, err
	}

	wallet, err := OpenWallet(cfg.WalletPath, cfg.LockAddr, cfg.RPCNode, cfg.WalletPass)
	if err != nil {
		return nil, err
	}

	listener := NewListener(ctx, client, bypass.PACTUS, highway, startBlock, cfg.LockAddr)
	bridge := NewBridge(wallet, b, bypass.PACTUS)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,
		bridge:   bridge,

		ctx: ctx,
	}, nil
}
