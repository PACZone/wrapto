package core

import (
	"context"
	"log"

	"github.com/PACZone/teleport/bridge"
	pactus "github.com/PACZone/teleport/client/pactus"
	polygon "github.com/PACZone/teleport/client/polygon"
	"github.com/PACZone/teleport/config"
	"github.com/PACZone/teleport/database"
	pactusListener "github.com/PACZone/teleport/listener/pactuslistener"
	polygonListener "github.com/PACZone/teleport/listener/polygonlistener"
	"github.com/PACZone/teleport/order"
	"github.com/PACZone/teleport/wallet"
)

const PolygonChainID = 80001 // TODO: make me multi chain!

type Core struct {
	orderCh  chan order.Order
	wallet   *wallet.Wallet
	db       *database.DB
	pactusL  *pactusListener.PactusListener
	polygonL *polygonListener.PolygonListener
	pactusC  *pactus.Mgr
	polygonC *polygon.Client
	bridge   *bridge.Bridge
}

func NewCore(pactusBlockStart, polygonOrderStart int, filepaths ...string) (*Core, error) {
	ctx := context.Background()

	cfg, err := config.LoadConfig(filepaths...)
	if err != nil {
		return nil, err
	}

	orderCh := make(chan order.Order, 10)

	w := wallet.Open(cfg.Wallet.Path, cfg.Wallet.Address,
		cfg.PacLsn.RPCURLS[0], cfg.Wallet.Password)

	db, err := database.NewDB(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	pactusCm := pactus.NewClientMgr(ctx)
	for _, c := range cfg.PacLsn.RPCURLS {
		pc, err := pactus.NewClient(c)
		if err != nil {
			log.Printf("can't connect to: %s", c)
		}
		pactusCm.AddClient(pc)
	}

	pactusL := pactusListener.NewPactusListener(pactusCm, orderCh,
		uint32(pactusBlockStart), cfg.PacLsn.BridgeAddress, db)

	polygonC, err := polygon.NewClient(cfg.PolLsn.RPCURL,
		cfg.PolLsn.PrivateKey, cfg.PolLsn.ContractAddress, PolygonChainID)
	if err != nil {
		return nil, err
	}

	polygonL := polygonListener.NewPolygonListener(uint32(polygonOrderStart), polygonC, orderCh, *db)

	brg := bridge.NewBridge(pactusCm, polygonC, orderCh, *w, *db)

	return &Core{
		orderCh:  orderCh,
		wallet:   w,
		db:       db,
		pactusL:  pactusL,
		polygonL: polygonL,
		pactusC:  pactusCm,
		polygonC: polygonC,
		bridge:   brg,
	}, nil
}

func (c *Core) Start() {
	go c.bridge.Start()
	go c.polygonL.Start()
	go c.pactusL.Start()
}

// TODO: implement STOP method for me!
