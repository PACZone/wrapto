package database_test

import (
	"os"
	"testing"

	"github.com/PACZone/wrapto/database"
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

func TestNewDB(t *testing.T) { // TODO: REMOVE ME LATER
	db := setup(t)
	assert.NotNil(t, db)
}
