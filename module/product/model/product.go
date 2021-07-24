package model

type Product struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Stock     int32  `json:"stock" db:"stock"`
	Price     int64  `json:"price" db:"price"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
