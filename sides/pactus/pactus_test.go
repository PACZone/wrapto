package pactus_test

import (
	"testing"

	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/stretchr/testify/assert"
)

func TestParseMemo(t *testing.T) {
	memos := []struct {
		memo string
		addr string
		dest bypass.Name
		err  bool
	}{
		{
			memo: "0x890189B78F2639a2A407C5F089481DB92A028608@POLYGON",
			addr: "0x890189B78F2639a2A407C5F089481DB92A028608",
			dest: bypass.POLYGON,
			err:  false,
		},
		{
			memo: "sergijwerpgij8",
			addr: "",
			dest: "",
			err:  true,
		},
		{
			memo: "a@2@r",
			addr: "",
			dest: "",
			err:  true,
		},
	}

	for _, m := range memos {
		d, err := pactus.ParseMemo(m.memo)
		if m.err {
			assert.Error(t, err)
			assert.Nil(t, d)

			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, d.Addr, m.addr)
		assert.Equal(t, d.BypassName, m.dest)
	}
}
