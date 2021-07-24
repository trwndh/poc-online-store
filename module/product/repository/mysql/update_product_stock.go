package mysql

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/product/model"
)

func (p *ProductRepo) UpdateProductStock(ctx context.Context, productID int64, stock int32) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.product.UpdateProductStock")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	var product model.Product
	err = tx.GetContext(ctx, &product, "SELECT * FROM product WHERE ID = ? FOR UPDATE", productID)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err
		}
		return err
	}

	query := "UPDATE product SET stock = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, stock, productID)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err
		}
		return err
	}

	_ = tx.Commit()
	return nil
}
