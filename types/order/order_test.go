package order_test

import (
	"testing"

	"github.com/PACZone/wrapto/types"
	"github.com/PACZone/wrapto/types/order"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	ord, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 10e9)
	assert.NoError(t, err)
	assert.Equal(t, order.CREATED, ord.Status)
}

func TestBasicCheck(t *testing.T) {
	t.Run("everything ok", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 10e9)
		assert.NoError(t, err)
	})

	t.Run("amount less than minium fee", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 0)
		assert.Error(t, err)
	})

	t.Run("amount equal to minium fee", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", types.MinimumFee)
		assert.Error(t, err)
	})
}

func TestFee(t *testing.T) {
	feeAndAmts := []struct { // better name?
		Fee    uint64
		Amount uint64
	}{
		{
			Amount: 1_903_076_060_983,
			Fee:    types.MaximumFee,
		},
		{
			Amount: 2_874_345_000,
			Fee:    types.MinimumFee,
		},
		{
			Amount: 200e9,
			Fee:    1e9,
		},
	}

	for _, fa := range feeAndAmts {
		ord, err := order.NewOrder("", "", "", fa.Amount)
		assert.NoError(t, err)

		assert.Equal(t, fa.Fee, ord.Fee())
	}
}

func TestAmount(t *testing.T) {
	feeAndAmts := []struct { // better name?
		Fee    uint64
		Amount uint64
	}{
		{
			Amount: 1_903_076_060_983,
			Fee:    types.MaximumFee,
		},
		{
			Amount: 2_874_345_000,
			Fee:    types.MinimumFee,
		},
		{
			Amount: 200e9,
			Fee:    1e9,
		},
	}

	for _, fa := range feeAndAmts {
		ord, err := order.NewOrder("", "", "", fa.Amount)
		assert.NoError(t, err)

		assert.Equal(t, (fa.Amount - fa.Fee), ord.Amount())
	}
}
