package mysql

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/product/model"
)

func (p *ProductRepo) AddProduct(ctx context.Context, input model.Product) (product model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.product.AddProduct")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return input, err
	}

	query := "INSERT INTO product(name, price, stock) VALUES(?,?,?)"
	res, err := tx.ExecContext(ctx, query, input.Name, input.Price, input.Stock)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return input, err
		}
		return input, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return input, err
		}
		return input, err
	}

	err = tx.GetContext(ctx, &product, "SELECT * FROM product WHERE id = ?", insertID)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return input, err
		}
		return input, err
	}

	_ = tx.Commit()
	return product, nil
}
