package dal

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

type TStruct struct {
	A int     `db:"a"`
	B string  `db:"b"`
	C bool    `db:"c"`
	D float32 `db:"d"`
}

func (t TStruct) Scan(values []any) Scaner {
	t.A = int(values[0].(int32))
	t.B = values[1].(string)
	t.C = values[2].(bool)
	t.D = float32(values[3].(float64))
	return &t
}

const initialValues = "VALUES (1, '2', true, 1.2), (3, '4' , false, 5.6)"
const testTable = "foo"
const createTableQ = "CREATE TABLE foo(a integer, b varchar(50), c bool, d double precision)"
const dropTableQ = "DROP TABLE foo"

func createTestTable(conn *pgx.Conn) {
	conn.Exec(context.Background(), createTableQ)
	conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s %s", testTable, initialValues))
}
func dropTestTable(conn *pgx.Conn) {
	conn.Exec(context.Background(), dropTableQ)
}

func TestPrepareQuery(t *testing.T) {
	d := NewDal[TStruct]("sometable")
	compiled, _, err := d.ReadDal.baseSelectQuery.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM sometable", compiled)
	t.Cleanup(func() {

	})
}

func TestFind(t *testing.T) {
	var ctx = context.TODO()
	d := NewDal[TStruct](testTable)
	createTestTable(d.ReadDal.conn)
	res, err := d.Find(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, 1, res[0].A)
	assert.Equal(t, "2", res[0].B)
	assert.Equal(t, true, res[0].C)
	assert.Equal(t, float32(1.2), res[0].D)
	t.Cleanup(func() {
		dropTestTable(d.ReadDal.conn)
	})
}

func TestFindOne(t *testing.T) {
	var ctx = context.TODO()
	d := NewDal[TStruct](testTable)
	createTestTable(d.ReadDal.conn)
	res, err := d.FindOne(ctx, "a = ?", 3)
	assert.NoError(t, err)
	assert.Equal(t, 3, res.A)
	assert.Equal(t, "4", res.B)
	assert.Equal(t, false, res.C)
	assert.Equal(t, float32(5.6), res.D)
	t.Cleanup(func() {
		dropTestTable(d.ReadDal.conn)
	})
}

func TestInsert(t *testing.T) {
	var ctx = context.TODO()
	d := NewDal[TStruct](testTable)
	createTestTable(d.ReadDal.conn)
	var structs = []TStruct{{
		A: -1,
		B: "-1",
		C: true,
		D: -1.0,
	}, {
		A: -2,
		B: "-2",
		C: false,
		D: -2.0,
	}, {
		A: -3,
		B: "-3",
		C: false,
		D: -3.0,
	}}
	err := d.Insert(ctx, structs...)
	assert.NoError(t, err)
	res, err := d.FindOne(ctx, "a = ?", -3)
	assert.NoError(t, err)
	assert.Equal(t, -3, res.A)
	assert.Equal(t, "-3", res.B)
	assert.Equal(t, false, res.C)
	assert.Equal(t, float32(-3), res.D)
	t.Cleanup(func() {
		dropTestTable(d.ReadDal.conn)
	})
}

func TestUpdate(t *testing.T) {
	var ctx = context.TODO()
	d := NewDal[TStruct](testTable)
	createTestTable(d.ReadDal.conn)
	var updateStruct = TStruct{-10, "-10", true, -10.0}
	err := d.Update(ctx, updateStruct, "a = ?", 1)
	assert.NoError(t, err)
	t.Cleanup(func() {
		dropTestTable(d.ReadDal.conn)
	})
	res, err := d.Find(ctx, "a = ?", -10)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, -10, res[0].A)
}

func TestDelete(t *testing.T) {
	var ctx = context.TODO()
	d := NewDal[TStruct](testTable)
	createTestTable(d.ReadDal.conn)
	var structs = []TStruct{{
		A: -100,
		B: "-100",
		C: true,
		D: -100.0,
	}}
	err := d.Insert(ctx, structs...)
	assert.NoError(t, err)

	err = d.Delete(ctx, "a = ?", -100)
	assert.NoError(t, err)
	res, err := d.Find(ctx, "a = ?", 100)
	assert.NoError(t, err)
	assert.Len(t, res, 0)
	t.Cleanup(func() {
		dropTestTable(d.ReadDal.conn)
	})
}
