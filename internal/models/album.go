package models

import (
	"github.com/shopspring/decimal"
)

type Album struct {
	tableName struct{}        `pg:"music.albums"`
	Id        int             `db:"id"`
	Title     string          `db:"title"`
	Artist    string          `db:"artist"`
	Price     decimal.Decimal `db:"price"`
}
