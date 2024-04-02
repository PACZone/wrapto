package bridge

import (
	"math/big"

	pactus "github.com/PACZone/teleport/client/pactus"
	"github.com/PACZone/teleport/client/polygon"
	"github.com/PACZone/teleport/database"
	"github.com/PACZone/teleport/order"
	"github.com/PACZone/teleport/wallet"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge struct {
	pactusClient *pactus.Mgr

	polygonClient *polygon.Client
	orderCh       chan order.Order
	wallet        wallet.Wallet
	db            database.DB
}

func NewBridge(pactusCmgr *pactus.Mgr, polygonC *polygon.Client,
	orderCh chan order.Order,
	w wallet.Wallet, db database.DB,
) *Bridge {
	return &Bridge{
		pactusClient:  pactusCmgr,
		polygonClient: polygonC,
		orderCh:       orderCh,
		wallet:        w,
		db:            db,
	}
}

func (b *Bridge) Start() {
	for { //nolint
		// TODO FIX LINT ISSUE
		select {
		case ord := <-b.orderCh:
			if err := b.db.UpdateOrderStatus(ord.ID, order.PENDING); err != nil {
				panic(err) // TODO: must be graceful shutdown
			}

			b.processOrder(&ord)
		}
	}
}

func (b *Bridge) processOrder(ord *order.Order) {
	if !ord.IsValid() {
		return
	}

	if ord.Type == order.PACTUS_POLYGON { //nolint
		// TODO EXPORT ME
		amountBigInt := new(big.Int).SetUint64(ord.Amount)
		hash, err := b.polygonClient.Mint(*amountBigInt, common.HexToAddress(ord.DestinationAddr))
		if err != nil {
			err := b.db.UpdateOrderProcessedHashAndReason(ord.ID, "", err.Error(), order.FAILED)
			if err != nil {
				panic(err) // TODO: must be graceful shutdown
			}

			return
		}

		if err = b.db.UpdateOrderProcessedHashAndReason(ord.ID, hash, "SUCCESSFUL", order.COMPLETE); err != nil {
			panic(err) // TODO: must be graceful shutdown
		}
	} else if ord.Type == order.POLYGON_PACTUS {
		hash, err := b.wallet.TransferTransaction(ord.DestinationAddr, ord.ID, int64(ord.Amount))
		if err != nil {
			if err = b.db.UpdateOrderProcessedHashAndReason(ord.ID, "", err.Error(), order.FAILED); err != nil {
				panic(err) // TODO: must be graceful shutdown
			}

			return
		}

		if err = b.db.UpdateOrderProcessedHashAndReason(ord.ID, hash, "SUCCESSFUL", order.COMPLETE); err != nil {
			panic(err) // TODO: must be graceful shutdown
		}

		panic("END") // TODO REMOVE ME
	}
}
