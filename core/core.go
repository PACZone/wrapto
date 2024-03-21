package core

import (
	"github.com/PacmanHQ/teleport/bridge"
	pactusClient "github.com/PacmanHQ/teleport/client/pactusclient"
	polygonClient "github.com/PacmanHQ/teleport/client/polygonclient"
	"github.com/PacmanHQ/teleport/config"
	"github.com/PacmanHQ/teleport/database"
	pactusListener "github.com/PacmanHQ/teleport/listener/pactuslistener"
	polygonListener "github.com/PacmanHQ/teleport/listener/polygonlistener"
	"github.com/PacmanHQ/teleport/order"
	"github.com/PacmanHQ/teleport/wallet"
)

const PolygonChainID = 80001 // TODO: make me multi chain!

type Core struct {
	orderCh  chan order.Order
	wallet   *wallet.Wallet
	db       *database.DB
	pactusL  *pactusListener.PactusListener
	polygonL *polygonListener.PolygonListener
	pactusC  *pactusClient.PactusClient
	polygonC *polygonClient.PolygonClient
	bridge   *bridge.Bridge
}

func NewCore(pactusBlockStart, polygonOrderStart int, filepaths ...string) *Core {
	cfg, err := config.LoadConfig(filepaths...)
	if err != nil {
		panic(err) // TODO: make me logger.panic...
	}

	orderCh := make(chan order.Order, 10)

	w := wallet.Open(cfg.Wallet.Path, cfg.Wallet.Address,
		cfg.PacLsn.RPCURLS[0], cfg.Wallet.Password)

	db, err := database.NewDB(cfg.DBPath)
	if err != nil {
		panic(err)
	}

	pactusC := pactusClient.NewPactusClient()
	for _, n := range cfg.PacLsn.RPCURLS {
		err = pactusC.AddClient(n)
		if err != nil {
			panic(err) // TODO: must be graceful shutdown
		}
	}

	pactusL := pactusListener.NewPactusListener(pactusC, orderCh,
		uint32(pactusBlockStart), cfg.PacLsn.BridgeAddress, db)

	polygonC, err := polygonClient.NewPolygonClient(cfg.PolLsn.RPCURL,
		cfg.PolLsn.PrivateKey, cfg.PolLsn.ContractAddress, PolygonChainID)
	if err != nil {
		panic(err)
	}

	polygonL := polygonListener.NewPolygonListener(uint32(polygonOrderStart), polygonC, orderCh, *db)

	brg := bridge.NewBridge(*pactusC, polygonC, orderCh, *w, *db)

	return &Core{
		orderCh:  orderCh,
		wallet:   w,
		db:       db,
		pactusL:  pactusL,
		polygonL: polygonL,
		pactusC:  pactusC,
		polygonC: polygonC,
		bridge:   brg,
	}
}

func (c *Core) Start() {
	go c.bridge.Start()
	go c.polygonL.Start()
	go c.pactusL.Start()
}

// TODO: implement STOP method for me!
