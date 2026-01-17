package album

import (
	_ "embed"

	"github.com/jmoiron/sqlx"
)

type AlbumRepository interface {
	Read() ([]Album, error)
}

type albumRepository struct {
	db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) AlbumRepository {
	return &albumRepository{db: db}
}

//go:embed queries/get_albums.sql
var getAlbumsQuery string

func (r *albumRepository) Read() ([]Album, error) {
	albums := []Album{}
	album := Album{}
	rows, _ := r.db.Queryx(getAlbumsQuery)
	for rows.Next() {
		err := rows.StructScan(&album)
		if err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}
	return albums, nil
}
