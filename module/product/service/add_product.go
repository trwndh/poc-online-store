package service

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/trwndh/poc-online-store/module/product/model"
)

func (s *service) AddProduct(ctx context.Context, input model.Product) (product model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.product.AddProduct")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	if input.Stock == 0 || input.Price == 0 || input.Name == "" {
		return product, errors.New("invalid input given")
	}

	product, err = s.productRepo.AddProduct(ctx, input)
	if err != nil {
		return product, err
	}

	return product, nil
}
