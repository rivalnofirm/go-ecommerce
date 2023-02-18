package product

import (
	"context"
	"product-service/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.Product, string, error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	Update(ctx context.Context, ar *models.Product) error
	Store(context.Context, *models.Product) error
	Delete(ctx context.Context, id int64) error
}
