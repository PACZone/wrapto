package pactus_test

import (
	"testing"

	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/stretchr/testify/assert"
)

func TestParseMemo(t *testing.T) {
	memos := []struct {
		memo        string
		addr        string
		dest        bypass.Name
		expectError bool
	}{
		{
			memo:        "0x682F5c4Bc85fEeC8D042D324960318553a47B24D@POLYGON",
			addr:        "0x682F5c4Bc85fEeC8D042D324960318553a47B24D",
			dest:        bypass.POLYGON,
			expectError: false,
		},
		{
			memo:        "0x682F5c4Bc85fEeC8D042D324960318553a47B24D@POLYGONN",
			addr:        "0x682F5c4Bc85fEeC8D042D324960318553a47B24D",
			dest:        bypass.POLYGON,
			expectError: true,
		},
		{
			memo:        "sergijwerpgij8",
			addr:        "",
			dest:        "",
			expectError: true,
		},
		{
			memo:        "a@2@r",
			addr:        "",
			dest:        "",
			expectError: true,
		},
	}

	for _, m := range memos {
		d, err := pactus.ParseMemo(m.memo)
		if m.expectError {
			assert.Error(t, err)
			assert.Empty(t, d)

			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, d.Addr, m.addr)
		assert.Equal(t, d.BypassName, m.dest)
	}
}
