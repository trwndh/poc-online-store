package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/trwndh/poc-online-store/module/product/model"
	"github.com/trwndh/poc-online-store/module/product/repository"
	logpkg "github.com/trwndh/poc-online-store/pkg/logger"
)

type service struct {
	ctx         context.Context
	productRepo repository.ProductRepository
}

type ProductService interface {
	GetProducts(ctx context.Context) (products []model.Product, err error)
	GetProductByID(ctx context.Context, productID int64) (product model.Product, err error)
	UpdateProductStock(ctx context.Context, productID int64, stock int32) (err error)
	AddProduct(ctx context.Context, input model.Product) (product model.Product, err error)
}

func NewProductService(ctx context.Context, productRepo repository.ProductRepository) ProductService {
	return &service{
		ctx:         ctx,
		productRepo: productRepo,
	}
}

var logger *zap.Logger = logpkg.Log
