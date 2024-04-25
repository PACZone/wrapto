package order_test

import (
	"testing"

	"github.com/PACZone/wrapto/types/order"
	"github.com/PACZone/wrapto/types/params"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	ord, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 10e9, order.PACTUS_POLYGON)
	assert.NoError(t, err)
	assert.Equal(t, order.PENDING, ord.Status)
}

func TestBasicCheck(t *testing.T) {
	t.Run("everything ok", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 10e9, order.PACTUS_POLYGON)
		assert.NoError(t, err)
	})

	t.Run("amount less than minium fee", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", 0, order.PACTUS_POLYGON)
		assert.Error(t, err)
	})

	t.Run("amount equal to minium fee", func(t *testing.T) {
		_, err := order.NewOrder("0x1234567", "pc1z123", "0xuoip", params.MinimumFee, order.PACTUS_POLYGON)
		assert.Error(t, err)
	})
}

func TestFee(t *testing.T) {
	feeAndAmounts := []struct { // better name?
		Fee    amount.Amount
		Amount amount.Amount
	}{
		{
			Amount: 1_903_076_060_983,
			Fee:    params.MaximumFee,
		},
		{
			Amount: 2_874_345_000,
			Fee:    params.MinimumFee,
		},
		{
			Amount: 200e9,
			Fee:    1e9,
		},
	}

	for _, fa := range feeAndAmounts {
		ord, err := order.NewOrder("", "", "", fa.Amount, order.PACTUS_POLYGON)
		assert.NoError(t, err)

		assert.Equal(t, fa.Fee, ord.Fee())
	}
}

func TestAmount(t *testing.T) {
	feeAndAmounts := []struct { // better name?
		Fee    amount.Amount
		Amount amount.Amount
	}{
		{
			Amount: 1_903_076_060_983,
			Fee:    params.MaximumFee,
		},
		{
			Amount: 2_874_345_000,
			Fee:    params.MinimumFee,
		},
		{
			Amount: 200e9,
			Fee:    1e9,
		},
	}

	for _, fa := range feeAndAmounts {
		ord, err := order.NewOrder("", "", "", fa.Amount, order.PACTUS_POLYGON)
		assert.NoError(t, err)

		assert.Equal(t, fa.Amount-fa.Fee, ord.Amount())
	}
}
