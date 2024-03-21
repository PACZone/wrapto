package bridge

import (
	"fmt"
	"math/big"

	pactusClient "github.com/PacmanHQ/teleport/client/pactus_client"
	polygonClient "github.com/PacmanHQ/teleport/client/polygon_client"
	"github.com/PacmanHQ/teleport/database"
	"github.com/PacmanHQ/teleport/order"
	"github.com/PacmanHQ/teleport/wallet"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge struct {
	pactusClient  pactusClient.PactusClient
	polygonClient polygonClient.PolygonClient
	orderCh       chan (order.Order)
	wallet        wallet.Wallet
	db            database.DB
}

func NewBridge(pactusC pactusClient.PactusClient, polygonC polygonClient.PolygonClient, orderCh chan (order.Order), w wallet.Wallet, db database.DB) *Bridge {
	return &Bridge{
		pactusClient:  pactusC,
		polygonClient: polygonC,
		orderCh:       orderCh,
		wallet:        w,
		db:            db,
	}
}

func (b *Bridge) Start() {
	fmt.Println("bridge start")
	for {
		select {
		case ord := <-b.orderCh:
			b.db.UpdateOrderStatus(ord.Id, order.PENDING)
			b.processOrder(ord)
		}
	}
}

func (b *Bridge) processOrder(ord order.Order) {
	if !ord.IsValid() {
		return
	}

	if ord.Type == order.PACTUS_POLYGON {
		amountBigInt := new(big.Int).SetUint64(ord.Amount)
		hash, err := b.polygonClient.Mint(*amountBigInt, common.HexToAddress(ord.DestinationAddr))
		if err != nil {
			b.db.UpdateOrderProcessedHashAndReason(ord.Id, "", err.Error(), order.FAILED)
			return
		}

		b.db.UpdateOrderProcessedHashAndReason(ord.Id, hash, "SUCCESSFUL", order.COMPLETE)

	} else if ord.Type == order.POLYGON_PACTUS {
		hash, err := b.wallet.TransferTransaction(ord.DestinationAddr, ord.Id, int64(ord.Amount))

		if err != nil {
			b.db.UpdateOrderProcessedHashAndReason(ord.Id, "", err.Error(), order.FAILED)
			return
		}

		b.db.UpdateOrderProcessedHashAndReason(ord.Id, hash, "SUCCESSFUL", order.COMPLETE)

		panic("END")
	}
}
