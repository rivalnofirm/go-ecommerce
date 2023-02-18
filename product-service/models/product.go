package models

import "time"

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"product_name" validate:"required"`
	Price     float32   `json:"price" validate:"required"`
	Brand     Brand     `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
