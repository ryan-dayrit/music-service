package orm

import (
	"github.com/go-pg/pg/v10"

	"music-service/internal/models"
)

type Repository interface {
	Read(id int) (*models.Album, error)
	Create(album models.Album) error
}

type repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Read(id int) (*models.Album, error) {
	album := &models.Album{Id: id}
	err := r.db.Model(album).WherePK().Select()
	return album, err
}

func (r *repository) Create(album models.Album) error {
	_, err := r.db.Model(&album).Insert()
	return err
}
