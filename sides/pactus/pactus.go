package pactus

import (
	"context"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
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
	b chan message.Message, env string, cfg config.PactusConfig,
	db *database.DB,
) (*Side, error) {
	if env == "dev" {
		crypto.AddressHRP = "tpc"
	}

	client, err := newClient(ctx, cfg.RPCNode)
	if err != nil {
		return nil, err
	}

	wallet, err := openWallet(cfg.WalletPath, cfg.LockAddr, cfg.RPCNode, cfg.WalletPass)
	if err != nil {
		return nil, err
	}

	listener := newListener(ctx, client, bypass.PACTUS, highway, startBlock, cfg.LockAddr, db)
	bridge := newBridge(ctx, wallet, b, bypass.PACTUS, db)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,
		bridge:   bridge,

		ctx: ctx,
	}, nil
}
