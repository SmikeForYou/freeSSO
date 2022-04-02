package dal

import "github.com/Masterminds/squirrel"

type Model interface {
	Scan(...any) Model
	Table() string
}

type SQLBuilder squirrel.StatementBuilderType
