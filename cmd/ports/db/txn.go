package db

import (
	"github.com/jackc/pgx/v5"
)

type QuerierTx interface {
	Querier
	WithTransaction(tx pgx.Tx) QuerierTx
}

func (q *Queries) WithTransaction(tx pgx.Tx) QuerierTx {
	return q.WithTx(tx)
}
