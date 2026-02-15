package orm

import (
	"github.com/go-pg/pg/v10"

	"music-service/internal/models"
)

type Repository interface {
	Create(album models.Album) error
	Read(id int) (*models.Album, error)
	Update(album models.Album) error
	Upsert(album models.Album) error
}

type repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(album models.Album) error {
	_, err := r.db.Model(&album).Insert()
	return err
}

func (r *repository) Read(id int) (*models.Album, error) {
	album := &models.Album{Id: id}
	err := r.db.Model(album).WherePK().Select()
	return album, err
}

func (r *repository) Update(album models.Album) error {
	_, err := r.db.Model(&album).WherePK().Update()
	return err
}

func (r *repository) Upsert(album models.Album) error {
	_, err := r.db.Model(&album).OnConflict("(id) DO UPDATE").Insert()
	return err
}
