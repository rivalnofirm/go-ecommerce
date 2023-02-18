package product

import (
	"context"
	"product-service/models"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Product, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetByName(ctx context.Context, title string) (*models.Product, error)
	Update(ctx context.Context, ar *models.Product) error
	Store(ctx context.Context, a *models.Product) error
	Delete(ctx context.Context, id int64) error
}
