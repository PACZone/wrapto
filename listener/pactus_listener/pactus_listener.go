package pactusListener

import (
	"context"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"

	pactusClient "github.com/PacmanHQ/teleport/client/pactus_client"
	"github.com/PacmanHQ/teleport/database"
	"github.com/PacmanHQ/teleport/order"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type PactusListener struct {
	client     *pactusClient.PactusClient
	orderCh   chan (order.Order)
	lastBlock  uint32
	bridgeAddr string
	ctx        context.Context
	db         *database.DB
}

func NewPactusListener(c *pactusClient.PactusClient, pactusCh chan(order.Order), lastBlock uint32, bridgeAddr string, db *database.DB) *PactusListener {
	return &PactusListener{
		client:     c,
		orderCh:   pactusCh,
		lastBlock:  lastBlock,
		bridgeAddr: bridgeAddr,
		ctx:        context.Background(),
		db:         db,
	}
}

func (p *PactusListener) Start() {
	for {
		p.processOrder()
	}
}

func (p *PactusListener) processOrder() {
	c, err := p.isRepeatedBlock(p.lastBlock)
	if err != nil || c {
		return
	}

	b, err := p.client.GetBlock(p.ctx, p.lastBlock, pactus.BlockVerbosity_BLOCK_TRANSACTIONS)
	if err != nil {
		return
	}

	extractedOrder := p.extractOrders(b.Txs)

	for _, order := range extractedOrder {
		err := p.db.AddOrder(order)
		if err != nil {
			//gracefull
		}
		p.orderCh <- order
	}

	a := p.db.CreateListened(0, int(p.lastBlock), len(extractedOrder))
	if a != nil {
		//gracefull
	}

	p.lastBlock++
}

func (p *PactusListener) isRepeatedBlock(block uint32) (bool, error) {
	lastBlockHeight, err := p.client.GetBlockchainHeight(p.ctx)
	if err != nil {
		return true, err
	}

	return block > lastBlockHeight, nil
}

func (p *PactusListener) extractOrders(txs []*pactus.TransactionInfo) []order.Order {

	var correctOrder []order.Order

	for _, tx := range txs {
		if tx.GetTransfer() != nil && tx.GetTransfer().Receiver == p.bridgeAddr {
			destAddr, err := detectDest(tx.Memo)
			if err != nil {
				//log
				continue
			}
			n, err := order.NewOrder(hex.EncodeToString(tx.Id), order.PACTUS_POLYGON, tx.GetTransfer().Sender, uint64(tx.GetTransfer().Amount), destAddr, p.lastBlock, p.bridgeAddr)
			if err != nil {
				//log
				continue
			}
			correctOrder = append(correctOrder, *n)
		}
	}

	return correctOrder
}

func detectDest(memo string) (string, error) {
	a := strings.Split(memo, ":")
	var addr string

	switch a[0] {
	case "POL":
		addr = a[1]
	default:
		return "", errors.New("invalid dest")
	}
	if isValidDestination(addr) {
		return a[1], nil
	} else {
		return "", errors.New("invalid dest")
	}
}

func isValidDestination(address string) bool {
	regex := regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	return regex.MatchString(address)
}
