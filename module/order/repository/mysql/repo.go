package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/trwndh/poc-online-store/module/order/repository"
)

type OrderRepo struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewOrderRepository(ctx context.Context, db *sqlx.DB) repository.OrderRepository {
	return &OrderRepo{
		ctx: ctx,
		db:  db,
	}
}
