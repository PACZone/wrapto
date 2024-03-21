package core

import (
	"os"

	"github.com/PacmanHQ/teleport/bridge"
	pactusClient "github.com/PacmanHQ/teleport/client/pactusclient"
	polygonClient "github.com/PacmanHQ/teleport/client/polygonclient"
	"github.com/PacmanHQ/teleport/database"
	pactusListener "github.com/PacmanHQ/teleport/listener/pactuslistener"
	polygonListener "github.com/PacmanHQ/teleport/listener/polygonlistener"
	"github.com/PacmanHQ/teleport/order"
	"github.com/PacmanHQ/teleport/wallet"
	"github.com/joho/godotenv"
)

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

var pactusNodes = []string{
	"bootstrap1.pactus.org:50051", "bootstrap2.pactus.org:50051",
	"bootstrap3.pactus.org:50051", "bootstrap4.pactus.org:50051",
	"151.115.110.114:50051", "188.121.116.247:50051",
}

func NewCore() *Core {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	orderCh := make(chan order.Order, 10)

	w := wallet.Open(os.Getenv("WALLET_PATH"), os.Getenv("WALLET_ADDRESS"),
		os.Getenv("PACTUS_NODE"), os.Getenv("WALLET_PASSWORD"))

	db, err := database.NewDB(os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}

	pactusC := pactusClient.NewPactusClient()
	for _, n := range pactusNodes {
		err = pactusC.AddClient(n)
		if err != nil {
			panic(err) // TODO: must be graceful shutdown
		}
	}

	pactusL := pactusListener.NewPactusListener(pactusC, orderCh,
		442320, os.Getenv("PACTUS_BRIDGE_ADDRESS"), db)

	polygonC, err := polygonClient.NewPolygonClient(os.Getenv("POLYGON_RPC"),
		os.Getenv("POLYGON_PRIVATE_KEY"), os.Getenv("POLYGON_CONTRACT_ADDRESS"), 80001)
	if err != nil {
		panic(err)
	}

	polygonL := polygonListener.NewPolygonListener(1, polygonC, orderCh, *db)

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
