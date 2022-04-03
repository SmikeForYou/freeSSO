package dal

import (
	"context"
	"freeSSO/internal/app/connections"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type ReadDal[T any] struct {
	conn            *pgx.Conn
	tablename       string
	baseSelectQuery squirrel.SelectBuilder
}

func NewReadDal[T any](tablename string) *ReadDal[T] {
	return &ReadDal[T]{
		tablename:       tablename,
		conn:            connections.GetDbConn(),
		baseSelectQuery: squirrel.Select("*").From(tablename).PlaceholderFormat(squirrel.Dollar),
	}
}

func (d *ReadDal[T]) Find(ctx context.Context, where string, args ...any) ([]*T, error) {
	compiled, val, err := d.baseSelectQuery.Where(where, args...).ToSql()
	if err != nil {
		return nil, err
	}
	res := make([]*T, 0)
	err = pgxscan.Select(context.Background(), d.conn, &res, compiled, val...)
	rows, err := d.conn.Query(ctx, "")
	rows.Scan()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *ReadDal[T]) FindOne(ctx context.Context, where string, args ...any) (*T, error) {
	compiled, val, err := d.baseSelectQuery.Where(where, args...).ToSql()
	if err != nil {
		return nil, err
	}
	var r T
	err = pgxscan.Get(context.TODO(), d.conn, &r, compiled, val...)
	return &r, err
}

type InsertDal[T any] struct {
	conn           *pgx.Conn
	tablename      string
	baseInserQuery squirrel.InsertBuilder
}

func NewInsertDal[T any](tablename string) *InsertDal[T] {
	return &InsertDal[T]{
		tablename:      tablename,
		conn:           connections.GetDbConn(),
		baseInserQuery: squirrel.Insert(tablename).PlaceholderFormat(squirrel.Dollar),
	}
}

func (d *InsertDal[T]) Insert(ctx context.Context, target ...T) error {
	insertQ := d.baseInserQuery
	for _, t := range target {
		vals, err := ToArr(t, "db")
		if err != nil {
			return err
		}
		insertQ = insertQ.Values(vals...)
	}
	q, vals, err := insertQ.ToSql()
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(ctx, q, vals...)
	return err
}

type UpdateDal[T any] struct {
	conn            *pgx.Conn
	tablename       string
	baseUpdateQuery squirrel.UpdateBuilder
}

func NewUpdateDal[T any](tablename string) *UpdateDal[T] {
	return &UpdateDal[T]{
		tablename:       tablename,
		conn:            connections.GetDbConn(),
		baseUpdateQuery: squirrel.Update(tablename).PlaceholderFormat(squirrel.Dollar),
	}
}

func (d *UpdateDal[T]) Update(ctx context.Context, target T, q string, vals ...any) error {
	valMap, err := ToMap(target, "db")
	if err != nil {
		return err
	}
	compiled, vls, err := d.baseUpdateQuery.SetMap(valMap).Where(q, vals...).ToSql()
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(ctx, compiled, vls...)
	return err
}

type DeleteDal struct {
	conn            *pgx.Conn
	tablename       string
	baseDeleteQuery squirrel.DeleteBuilder
}

func NewDeleteDal(tablename string) *DeleteDal {
	return &DeleteDal{
		tablename:       tablename,
		conn:            connections.GetDbConn(),
		baseDeleteQuery: squirrel.Delete(tablename).PlaceholderFormat(squirrel.Dollar),
	}
}

func (d *DeleteDal) Delete(ctx context.Context, q string, vals ...any) error {
	compiled, vls, err := d.baseDeleteQuery.Where(q, vals...).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(ctx, compiled, vls...)
	return err
}

type Dal[T any] struct {
	*ReadDal[T]
	*InsertDal[T]
	*UpdateDal[T]
	*DeleteDal
}

func NewDal[T any](tablename string) *Dal[T] {
	return &Dal[T]{
		ReadDal:   NewReadDal[T](tablename),
		InsertDal: NewInsertDal[T](tablename),
		UpdateDal: NewUpdateDal[T](tablename),
		DeleteDal: NewDeleteDal(tablename),
	}
}
