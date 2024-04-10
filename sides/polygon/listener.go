package polygon

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/PACZone/wrapto/types/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pactus-project/pactus/types/amount"
)

type Listener struct {
	client          *Client
	db              *database.DB
	bypassName      bypass.Name
	nextOrderNumber uint32
	highway         chan message.Message

	ctx context.Context
}

func newListener(ctx context.Context,
	client *Client, bp bypass.Name, highway chan message.Message, startOrder uint32, db *database.DB,
) *Listener {
	return &Listener{
		client:          client,
		db:              db,
		bypassName:      bp,
		nextOrderNumber: startOrder,
		highway:         highway,
		ctx:             ctx,
	}
}

func (l *Listener) Start() error {
	logger.Info("starting listener", "actor", l.bypassName)

	for {
		select {
		case <-l.ctx.Done():
			logger.Info("stopping listener", "actor", l.bypassName, "nextOrder", l.nextOrderNumber)

			return nil
		default:
			if err := l.processOrder(); err != nil {
				logger.Error("can't process block on listener", "actor", l.bypassName, "err", err)

				return err
			}
		}
	}
}

func (l *Listener) processOrder() error {
	o, err := l.client.Get(*big.NewInt(int64(l.nextOrderNumber)))
	if err != nil {
		return err
	}

	if o.Sender == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		time.Sleep(3 * time.Second)

		return nil
	}

	l.nextOrderNumber++

	id := strconv.FormatUint(uint64(l.nextOrderNumber), 10)

	if exist, err := l.checkOrderExist(id); err != nil {
		return err
	} else if exist {
		logger.Warn("error repetitive transaction", "actor", l.bypassName, "txHash", id)

		return nil
	}

	logger.Info("processing new message on listener", "actor", l.bypassName, "orderNumber", l.nextOrderNumber)

	amt := o.Amount.Int64()
	sender := o.Sender.Hex()
	ord, err := order.NewOrder(id, sender, o.DestinationAddress, amount.Amount(amt))
	if err != nil {
		dbErr := l.db.UpdateOrderStatus(ord.ID, order.FAILED)
		if dbErr != nil {
			return dbErr
		}

		dbErr = l.db.AddLog("", string(l.bypassName), fmt.Sprintf("failed to create order: %s", id), err.Error())
		if dbErr != nil {
			return dbErr
		}

		return nil
	}

	id, err = l.db.AddOrder(ord)
	if err != nil {
		return err
	}

	err = l.db.AddLog(id, string(l.bypassName), "order created", "")
	if err != nil {
		return err
	}

	msg := message.NewMessage(params.MainBypass, l.bypassName, ord)
	l.highway <- msg

	logger.Info("new message passed to pactus", "actor", l.bypassName, "orderID", ord.ID)

	err = l.db.AddLog(id, string(l.bypassName), "sent order to highway", "")
	if err != nil {
		return err
	}

	err = l.db.UpdatePolygonState(l.nextOrderNumber - 1)
	if err != nil {
		return err
	}

	return nil
}

func (l *Listener) checkOrderExist(id string) (bool, error) {
	isExist, err := l.db.IsOrderExist(id)
	if err != nil {
		return false, err
	}

	return isExist, nil
}
