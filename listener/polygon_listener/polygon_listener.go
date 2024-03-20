package polygonListener

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/PacmanHQ/teleport/client/polygon_client"
	"github.com/PacmanHQ/teleport/order"
	"github.com/ethereum/go-ethereum/common"
)

type PolygonListener struct {
	client    polygonClient.PolygonClient
	polygonCh chan (order.Order)
	lastOrder uint32
}

func NewPolygonListener(startFrom uint32, c polygonClient.PolygonClient, polygonCh *chan (order.Order)) *PolygonListener {
	return &PolygonListener{
		lastOrder: startFrom,
		client:    c,
		polygonCh: *polygonCh,
	}
}

func (p *PolygonListener) Start() {
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		go func() {
			<- time.After(5*time.Second)
			ord, err := p.client.GetOrder(*big.NewInt(int64(p.lastOrder)))
			if err != nil || ord.Sender == common.HexToAddress("0x0000000000000000000000000000000000000000") {
				return
			}
			fmt.Println(p.lastOrder)

			o, err := order.NewOrder("", order.POLYGON_PACTUS, ord.Sender.String(), ord.Amount.Uint64(), ord.DestinationAddress, p.lastOrder)
			if err != nil {
				return
			}
			fmt.Println(ord.Sender.String())
			fmt.Println(ord.Amount.Uint64())
			p.polygonCh <- *o

			p.lastOrder++
			wg.Done()
		}()
		wg.Wait()
	}
}
