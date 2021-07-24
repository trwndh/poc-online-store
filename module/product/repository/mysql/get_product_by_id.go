package mysql

import (
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/product/model"
)

func (p *ProductRepo) GetProductByID(ctx context.Context, productID int64) (product model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.product.GetProductByID")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return product, err
	}
	err = tx.GetContext(ctx, &product, "SELECT * FROM product WHERE id = ? FOR UPDATE", productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		if err2 := tx.Rollback(); err2 != nil {
			return product, err
		}
		return product, err
	}

	_ = tx.Commit()
	return product, nil
}
