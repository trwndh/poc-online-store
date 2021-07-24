package service

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (s *service) UpdateProductStock(ctx context.Context, productID int64, stock int32) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.product.UpdateProductStock")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	if productID == 0 || stock == 0 {
		return errors.New("invalid input given")
	}

	err = s.productRepo.UpdateProductStock(ctx, productID, stock)
	if err != nil {
		return err
	}

	return nil
}
