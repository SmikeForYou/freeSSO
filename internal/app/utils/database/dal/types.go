package dal

import "context"

type Scaner interface {
	Scan([]any) Scaner
}

type Reader[T any] interface {
	Find(context.Context, string, ...any) ([]*T, error)
	FindOne(context.Context, string, ...any) (*T, error)
}

type Inserter[T any] interface {
	Insert(context.Context, ...T) error
}

type Updater[T any] interface {
	Update(context.Context, T, string, ...any) error
}

type Deleter interface {
	Delete(context.Context, string, ...any) error
}
