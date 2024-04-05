package polygon

import (
	"context"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

const (
	mainChainID   int64 = 137
	mumbaiChainID int64 = 80001
)

type Side struct {
	client   *Client
	listener *Listener
	bridge   Bridge
	db       *database.DB
	highway  chan message.Message
	bypass   chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context, highway chan message.Message, startOrder uint32,
	bp chan message.Message, env string, cfg config.PolygonConfig,db *database.DB,
) (*Side, error) {
	chainID := mainChainID
	if env == "dev" {
		chainID = mumbaiChainID
	}

	client, err := newClient(cfg.RPCNode, cfg.PrivateKey, cfg.ContractAddr, chainID)
	if err != nil {
		return nil, err
	}

	listener := newListener(ctx, client, bypass.POLYGON, highway, startOrder)
	bridge := newBridge(bp, bypass.POLYGON, client)

	return &Side{
		client:   client,
		listener: listener,
		bridge:   bridge,
		highway:  highway,
		bypass:   bp,

		ctx: ctx,
	}, nil
}
