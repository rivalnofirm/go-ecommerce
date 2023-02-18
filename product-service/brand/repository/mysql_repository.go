package repository

import (
	"context"
	"database/sql"

	"product-service/brand"
	"product-service/models"

	"github.com/sirupsen/logrus"
)

type mysqlBrandRepo struct {
	DB *sql.DB
}

func NewMysqlBrandRepository(db *sql.DB) brand.Repository {
	return &mysqlBrandRepo{
		DB: db,
	}
}

func (m *mysqlBrandRepo) getOne(ctx context.Context, query string, args ...interface{}) (*models.Brand, error) {

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	a := &models.Brand{}

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return a, nil
}

func (m *mysqlBrandRepo) GetByID(ctx context.Context, id int64) (*models.Brand, error) {
	query := `SELECT id, name, created_at, updated_at FROM brand WHERE id=?`
	return m.getOne(ctx, query, id)
}
