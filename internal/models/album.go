package models

import (
	"github.com/shopspring/decimal"
)

type Album struct {
	Id     int             `db:"id"`
	Title  string          `db:"title"`
	Artist string          `db:"artist"`
	Price  decimal.Decimal `db:"price"`
}
