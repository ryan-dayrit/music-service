package models

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Album struct {
	tableName struct{}        `pg:"music.albums"`
	Id        int             `db:"id"`
	Title     string          `db:"title"`
	Artist    string          `db:"artist"`
	Price     decimal.Decimal `db:"price"`
}

func (a *Album) String() string {
	return fmt.Sprintf("Album{Id: %d, Title: %s, Artist: %s, Price: %s}", a.Id, a.Title, a.Artist, a.Price.String())
}
