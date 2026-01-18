package album

import (
	_ "embed"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Read() ([]Album, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

//go:embed queries/get_albums.sql
var getAlbumsQuery string

func (r *repository) Read() ([]Album, error) {
	albums := []Album{}
	rows, _ := r.db.Queryx(getAlbumsQuery)
	for rows.Next() {
		album := Album{}
		err := rows.StructScan(&album)
		if err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}
	return albums, nil
}
