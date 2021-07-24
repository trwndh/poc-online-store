package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/trwndh/poc-online-store/module/product/model"
)

func (s *service) GetProducts(ctx context.Context) (products []model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.product.GetProducts")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	products, err = s.productRepo.GetProducts(ctx)
	if err != nil {
		logger.Error("error when get products: " + err.Error())
		return products, err
	}

	return products, nil
}
