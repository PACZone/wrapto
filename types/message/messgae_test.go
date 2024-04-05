package message_test

import (
	"testing"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) message.Message {
	t.Helper()

	m := message.NewMessage(bypass.PACTUS, bypass.POLYGON, &order.Order{})
	require.NotNil(t, m)

	return m
}

func TestValidate(t *testing.T) {
	m := setup(t)

	t.Run("invalid destination", func(t *testing.T) {
		err := m.Validate(bypass.POLYGON)
		assert.Error(t, err)
	})

	t.Run("valid destination", func(t *testing.T) {
		err := m.Validate(bypass.PACTUS)
		assert.NoError(t, err)
	})
}
