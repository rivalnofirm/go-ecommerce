package models

type Brand struct {
	ID        int64  `json:"id"`
	Name      string `json:"brand_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
