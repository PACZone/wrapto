package main

import (
	"os"

	"github.com/PacmanHQ/teleport/bridge"
	pactusClient "github.com/PacmanHQ/teleport/client/pactus_client"
	polygonClient "github.com/PacmanHQ/teleport/client/polygon_client"
	"github.com/PacmanHQ/teleport/database"
	pactusListener "github.com/PacmanHQ/teleport/listener/pactus_listener"
	polygonListener "github.com/PacmanHQ/teleport/listener/polygon_listener"
	"github.com/PacmanHQ/teleport/order"
	"github.com/PacmanHQ/teleport/wallet"
	"github.com/joho/godotenv"
)

func main() {
	e := make(chan (string))
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	c := make(chan order.Order, 10)

	wallet := wallet.Open(os.Getenv("WALLET_PATH"), os.Getenv("WALLET_ADDRESS"), os.Getenv("PACTUS_NODE"), os.Getenv("WALLET_PASSWORD"))

	db, err := database.NewDB("./teleq.sqlite")
	if err != nil {
		panic(err)
	}

	p := pactusClient.NewPactusClient()
	p.AddClient(os.Getenv("PACTUS_NODE"))

	a, err := polygonClient.NewPolygonClient(os.Getenv("POLYGON_RPC"), os.Getenv("POLYGON_PRIVATE_KEY"), os.Getenv("POLYGON_CONTRACT_ADDRESS"), 80001)
	if err != nil {
		panic(err)
	}

	bridge := bridge.NewBridge(*p, *a, c, *wallet, *db)
	go bridge.Start()

	b := polygonListener.NewPolygonListener(1, *a, &c, *db)

	go b.Start()

	bb := pactusListener.NewPactusListener(*p, c, 442320, os.Getenv("PACTUS_BRIDGE_ADDRESS"), *db)

	go bb.Start()

	<-e
}
