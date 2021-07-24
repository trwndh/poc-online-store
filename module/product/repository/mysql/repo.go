package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/trwndh/poc-online-store/module/product/repository"
)

type ProductRepo struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewProductRepository(ctx context.Context, db *sqlx.DB) repository.ProductRepository {
	return &ProductRepo{
		ctx: ctx,
		db:  db,
	}
}
