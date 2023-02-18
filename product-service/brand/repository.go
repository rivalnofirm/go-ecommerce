package brand

import (
	"context"
	"product-service/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Brand, error)
}
