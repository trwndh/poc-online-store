package mysql

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/product/model"
)

func (p *ProductRepo) GetProducts(ctx context.Context) (products []model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.product.GetProducts")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return products, err
	}
	pp := []model.Product{}
	err = tx.SelectContext(ctx, &pp, "SELECT * FROM product")
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return products, err
		}
		return products, err
	}

	_ = tx.Commit()
	return pp, nil
}
