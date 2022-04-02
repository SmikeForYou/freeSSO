package dal_test

import (
	dal "freeSSO/internal/app/utils/database/dal"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

type TStruct struct {
	A int
	B string
	C bool
	D float32
}

func TestFind(t *testing.T) {
	d := dal.NewDal[TStruct]("(SELECT 1,'2',true, 1.2) as bar(A, B, C, D)")
	q := squirrel.Select("*")
	res, err := d.Find(q)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 1, res[0].A)
	assert.Equal(t, "2", res[0].B)
	assert.Equal(t, true, res[0].C)
	assert.Equal(t, float32(1.2), res[0].D)
}
