package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	pgxLogger := GetPgxLogger()
	assert.NotNil(t, pgxLogger)
	logg := GetLogger()
	assert.NotNil(t, logg)
}
