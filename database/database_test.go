package database_test

import (
	"os"
	"testing"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *database.DB {
	t.Helper()

	file, err := os.CreateTemp("", "temp-db")
	require.NoError(t, err)

	db, err := database.NewDB(file.Name())
	require.NoError(t, err)

	return db
}

func TestAddOrder(t *testing.T) {
	db := setup(t)

	ordID, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	assert.NoError(t, err)

	o, err := db.AddOrder(ordID)
	assert.NoError(t, err)

	assert.Equal(t, ordID.ID, o)
}

func TestAddLog(t *testing.T) {
	db := setup(t)

	err := db.AddLog("dnslkn","POLYGON","this is abcd","")
	assert.NoError(t, err)
}

func TestAddLogForOrder(t *testing.T) {
	db := setup(t)

	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	assert.NoError(t, err)

	o, err := db.AddOrder(ord)
	assert.NoError(t, err)

	err = db.AddLog(o,"POLYGON","descriptivjerijw","trace")
	assert.NoError(t, err)
}

func TestUpdateOrderStatus(t *testing.T) {
	db := setup(t)

	newOrd, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	require.NoError(t, err)

	ordID, err := db.AddOrder(newOrd)
	require.NoError(t, err)

	err = db.UpdateOrderStatus(ordID, order.COMPLETE)
	require.NoError(t, err)

	updatedOrd, err := db.GetOrder(ordID)
	require.NoError(t, err)

	assert.Equal(t, order.COMPLETE, updatedOrd.Status)
}

func TestGetOrder(t *testing.T) {
	db := setup(t)

	newOrd, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	require.NoError(t, err)

	o, err := db.AddOrder(newOrd)
	require.NoError(t, err)

	retOrd, err := db.GetOrder(o)
	require.NoError(t, err)
	assert.Equal(t, retOrd.ID, newOrd.ID)
	assert.Equal(t, retOrd.TxHash, newOrd.TxHash)
	assert.Equal(t, retOrd.Amount, newOrd.OriginalAmount())
	assert.Equal(t, retOrd.Fee, newOrd.Fee())
	assert.Equal(t, retOrd.Sender, newOrd.Sender)
	assert.Equal(t, retOrd.Receiver, newOrd.Receiver)
}

func TestGetOrderWithLogs(t *testing.T) {
	db := setup(t)

	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	assert.NoError(t, err)

	ordID, err := db.AddOrder(ord)
	assert.NoError(t, err)

	err = db.AddLog(ordID,"POLYGON","descriptivjerijw","trace")

	assert.NoError(t, err)

	ordWithLogs, err := db.GetOrderWithLogs(ordID)
	assert.NoError(t, err)

	assert.Equal(t, len(ordWithLogs.Logs), 1)
}

func TestGetOrderLogs(t *testing.T) {
	db := setup(t)

	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
	assert.NoError(t, err)

	ordID, err := db.AddOrder(ord)
	assert.NoError(t, err)

	err = db.AddLog(ordID,"POLYGON","abcd","traceAbcd")

	assert.NoError(t, err)

	logs, err := db.GetOrderLogs(ordID)
	assert.NoError(t, err)

	l := logs[0]
	assert.Equal(t, l.Actor, "POLYGON")
	assert.Equal(t, l.Description, "abcd")
	assert.Equal(t, l.Trace, "traceAbcd")
	assert.Equal(t, l.OrderID, ordID)
}
