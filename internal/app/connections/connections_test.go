package connections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDbConnection(t *testing.T) {
	dbconn := GetDbConn()
	assert.NotNil(t, dbconn)
	assert.NotNil(t, pool)
}
