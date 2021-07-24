package repository

import (
	"context"

	"github.com/trwndh/poc-online-store/module/product/model"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) (products []model.Product, err error)
	GetProductByID(ctx context.Context, productID int64) (product model.Product, err error)
	UpdateProductStock(ctx context.Context, productID int64, stock int32) (err error)
	AddProduct(ctx context.Context, input model.Product) (product model.Product, err error)
}
