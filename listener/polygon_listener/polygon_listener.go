package polygonListener

import (
	"math/big"

	"github.com/PacmanHQ/teleport/client/polygon_client"
	"github.com/PacmanHQ/teleport/database"
	"github.com/PacmanHQ/teleport/order"
	"github.com/ethereum/go-ethereum/common"
)

type PolygonListener struct {
	client    polygonClient.PolygonClient
	polygonCh chan (order.Order)
	lastOrder uint32
	DB        database.DB
}

func NewPolygonListener(startFrom uint32, c polygonClient.PolygonClient, polygonCh *chan (order.Order), db database.DB) *PolygonListener {
	return &PolygonListener{
		lastOrder: startFrom,
		client:    c,
		polygonCh: *polygonCh,
		DB:        db,
	}
}

func (p *PolygonListener) Start() {
	for {
		p.processOrder()
	}
}

func (p *PolygonListener) processOrder() {
	ord, err := p.client.GetOrder(*big.NewInt(int64(p.lastOrder)))
	if err != nil || ord.Sender == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return
	}

	o, err := order.NewOrder("", order.POLYGON_PACTUS, ord.Sender.String(), ord.Amount.Uint64(), ord.DestinationAddress, p.lastOrder, "")
	if err != nil {
		return
	}

	a := p.DB.AddOrder(*o)
	if a != nil {
		//
	}

	p.polygonCh <- *o

	b := p.DB.CreateListened(1, int(p.lastOrder), 1)
	if b != nil {
		//
	}

	p.lastOrder++
}
