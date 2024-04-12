package database_test

// import (
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/PACZone/wrapto/database"
// 	"github.com/PACZone/wrapto/types/order"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/gorm"
// )

// func setup(t *testing.T) *database.DB {
// 	t.Helper()

// 	file, err := os.CreateTemp("", "temp-db")
// 	require.NoError(t, err)

// 	db, err := database.NewDB(file.Name())
// 	require.NoError(t, err)

// 	return db
// }

// func TestAddOrder(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	assert.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	assert.NoError(t, err)

// 	assert.Equal(t, ord.ID, ordID)
// }

// func TestAddLog(t *testing.T) {
// 	db := setup(t)

// 	err := db.AddLog("1", "POLYGON", "this is test desc", "trace")
// 	assert.NoError(t, err)
// }

// func TestAddLogForOrder(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	assert.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	assert.NoError(t, err)

// 	err = db.AddLog(ordID, "POLYGON", "desc", "trace")
// 	assert.NoError(t, err)
// }

// func TestUpdateOrderStatus(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("0xFFFF", "sender", "receiver", 2e9)
// 	require.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	require.NoError(t, err)

// 	err = db.UpdateOrderStatus(ordID, order.COMPLETE)
// 	require.NoError(t, err)

// 	updatedOrd, err := db.GetOrder(ordID)
// 	require.NoError(t, err)

// 	assert.Equal(t, order.COMPLETE, updatedOrd.Status)
// }

// func TestGetOrder(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	require.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	require.NoError(t, err)

// 	retrievedOrd, err := db.GetOrder(ordID)
// 	require.NoError(t, err)
// 	assert.Equal(t, retrievedOrd.ID, ord.ID)
// 	assert.Equal(t, retrievedOrd.TxHash, ord.TxHash)
// 	assert.Equal(t, retrievedOrd.Amount, ord.OriginalAmount())
// 	assert.Equal(t, retrievedOrd.Fee, ord.Fee())
// 	assert.Equal(t, retrievedOrd.Sender, ord.Sender)
// 	assert.Equal(t, retrievedOrd.Receiver, ord.Receiver)
// }

// func TestGetOrderWithLogs(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	assert.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	assert.NoError(t, err)

// 	err = db.AddLog(ordID, "POLYGON", "descriptivjerijw", "trace")

// 	assert.NoError(t, err)

// 	ordWithLogs, err := db.GetOrderWithLogs(ordID)
// 	assert.NoError(t, err)

// 	assert.Equal(t, len(ordWithLogs.Logs), 1)
// }

// func TestGetOrderLogs(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	assert.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	assert.NoError(t, err)

// 	err = db.AddLog(ordID, "POLYGON", "abcd", "traceAbcd")

// 	assert.NoError(t, err)

// 	logs, err := db.GetOrderLogs(ordID)
// 	assert.NoError(t, err)

// 	l := logs[0]
// 	assert.Equal(t, l.Actor, "POLYGON")
// 	assert.Equal(t, l.Description, "abcd")
// 	assert.Equal(t, l.Trace, "traceAbcd")
// 	assert.Equal(t, l.OrderID, ordID)
// }

// func TestIsOrderExist(t *testing.T) {
// 	db := setup(t)

// 	ord, err := order.NewOrder("aaa", "sendet", "rec", 20e9)
// 	assert.NoError(t, err)

// 	ordID, err := db.AddOrder(ord)
// 	assert.NoError(t, err)

// 	isExist, err := db.IsOrderExist("bbb")
// 	assert.Equal(t, isExist, false)
// 	assert.Error(t, gorm.ErrRecordNotFound, err)

// 	isExist, err = db.IsOrderExist("aaa")
// 	assert.Equal(t, isExist, true)
// 	assert.NoError(t, err)

// 	assert.Equal(t, ord.ID, ordID)
// }

// func TestUpdatePactusState(t *testing.T) {
// 	db := setup(t)

// 	state := database.State{
// 		Pactus:  0,
// 		Polygon: 0,
// 	}
// 	assert.NoError(t, db.Create(&state).Error)

// 	err := db.UpdatePactusState(42)
// 	assert.NoError(t, err)

// 	var updatedState database.State
// 	assert.NoError(t, db.First(&updatedState).Error)
// 	assert.Equal(t, uint32(42), updatedState.Pactus)
// }

// func TestUpdatePolygonState(t *testing.T) {
// 	db := setup(t)

// 	state := database.State{
// 		Pactus:  0,
// 		Polygon: 0,
// 	}
// 	assert.NoError(t, db.Create(&state).Error)

// 	err := db.UpdatePolygonState(100)
// 	assert.NoError(t, err)

// 	var updatedState database.State
// 	assert.NoError(t, db.First(&updatedState).Error)
// 	assert.Equal(t, uint32(100), updatedState.Polygon)
// }

// func TestGetState(t *testing.T) {
// 	db := setup(t)

// 	err := db.UpdatePolygonState(100)
// 	assert.NoError(t, err)

// 	err = db.UpdatePactusState(42)
// 	assert.NoError(t, err)

// 	retrievedState, err := db.GetState()
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint32(42), retrievedState.Pactus)
// 	assert.Equal(t, uint32(100), retrievedState.Polygon)
// }

// func TestGetLatestOrders(t *testing.T) {
// 	db := setup(t)

// 	for i := 0; i < 15; i++ {
// 		ord, err := order.NewOrder(fmt.Sprintf("tx%d", i), fmt.Sprintf("sender%d", i), fmt.Sprintf("receiver%d", i), 20e9)
// 		require.NoError(t, err)
// 		_, err = db.AddOrder(ord)
// 		require.NoError(t, err)
// 	}

// 	orders, err := db.GetLatestOrders(10)
// 	require.NoError(t, err)

// 	require.Len(t, orders, 10)

// 	for i := 0; i < 10; i++ {
// 		assert.Equal(t, fmt.Sprintf("tx%d", 14-i), orders[i].TxHash)
// 	}
// }

// func TestSearchOrders(t *testing.T) {
// 	db := setup(t)

// 	for i := 0; i < 5; i++ {
// 		ord, err := order.NewOrder(fmt.Sprintf("tx%d", i), fmt.Sprintf("sender%d", i), fmt.Sprintf("receiver%d", i), 20e9)
// 		require.NoError(t, err)
// 		_, err = db.AddOrder(ord)
// 		require.NoError(t, err)
// 	}

// 	orders, err := db.SearchOrders("tx2")
// 	require.NoError(t, err)

// 	require.Len(t, orders, 1)
// 	assert.Equal(t, "tx2", orders[0].TxHash)

// 	orders, err = db.SearchOrders("receiver3")
// 	require.NoError(t, err)

// 	require.Len(t, orders, 1)
// 	assert.Equal(t, "receiver3", orders[0].Receiver)
// }
