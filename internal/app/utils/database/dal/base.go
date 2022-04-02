package dal

import (
	"context"
	"freeSSO/internal/app/connections"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Dal[T any] struct {
	conn      *pgx.Conn
	tablename string
}

func (d *Dal[T]) prepareQuery(q squirrel.SelectBuilder) (string, []interface{}, error) {
	query := q.From(d.tablename).PlaceholderFormat(squirrel.Dollar)
	return query.ToSql()

}

func (d *Dal[T]) Find(q squirrel.SelectBuilder) ([]*T, error) {
	compiled, val, err := d.prepareQuery(q)
	if err != nil {
		return nil, err
	}
	res := make([]*T, 0)
	err = pgxscan.Select(context.Background(), d.conn, &res, compiled, val...)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func NewDal[T any](tablename string) *Dal[T] {
	return &Dal[T]{conn: connections.GetDbConnPool(), tablename: tablename}
}
