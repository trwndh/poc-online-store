package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/trwndh/poc-online-store/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/order/model"
)

func (o *OrderRepo) CheckoutNoLocking(ctx context.Context, userID int64, cartID int64, cartItems []model.CartItem) (registerPaymentID int64, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.order.checkout")
	defer span.Finish()
	logger := logger.Log
	// key map pair for id product and its requested quantity
	mapProductStock := make(map[int64]int32, 0)

	// map of product ids, to lock product row
	var productIDs []int64
	for _, product := range cartItems {
		mapProductStock[product.ProductID] = product.Quantity
		productIDs = append(productIDs, product.ProductID)
	}

	// checkout step:
	// 0. start transaction
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
			_ = tx.Rollback()
		}
	}()

	// 1. get product data by its ID and lock its row
	products, err := findProductNoLocking(ctx, tx, productIDs)
	if err != nil {
		logger.Error("error finding products " + err.Error())
		return 0, err
	}

	// 2. if quantity avail >= quantity requested, continue, else, return error with stock unavailable.
	errCh := make(chan error, 0)
	var wg sync.WaitGroup
	for _, product := range products {
		if mapProductStock[product.ID] <= product.Stock {
			// concurrent update
			wg.Add(1)
			go func() {
				defer wg.Done()

				q := "select stock from product where id = ?"
				var stockBefore, stockAfter int
				_ = tx.QueryRowContext(ctx, q, product.ID).Scan(&stockBefore)
				err := updateProductStock(ctx, tx, product.ID, mapProductStock[product.ID])
				if err != nil {
					errCh <- err
				}
				_ = tx.QueryRowContext(ctx, q, product.ID).Scan(&stockAfter)
				logger.Info(fmt.Sprintf("updating stock for %s (requested %d) => stock before : %d, stock after : %d", product.Name, mapProductStock[product.ID], stockBefore, stockAfter))
			}()
		} else {
			errMsg := fmt.Sprintf("insufficient stock for product %s, want %d but only available %d", product.Name, mapProductStock[product.ID], product.Stock)
			return 0, errors.New(errMsg)
		}

	}
	wg.Wait()
	// wait until all update done
	for i := 0; i < len(errCh); i++ {
		select {
		case err := <-errCh:
			if err != nil {
				logger.Error("error when updating stock :" + err.Error())
				return 0, err
			}
		}
	}

	// 4. insert checkout data to register_payment table. to give user time limit to pay
	jsonCartItem, _ := json.Marshal(cartItems)
	registerPaymentID, err = registerPayment(ctx, tx, cartID, userID, string(jsonCartItem))
	if err != nil {
		return 0, err
	}
	// 5. commit.
	tx.Commit()

	return registerPaymentID, nil
}

func findProductNoLocking(ctx context.Context, tx *sqlx.Tx, productIDs []int64) (products []model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.order.checkout.findProduct")
	defer span.Finish()

	query, args, err := sqlx.In("SELECT id, name, stock FROM product WHERE id IN (?)", productIDs)
	if err != nil {
		return []model.Product{}, err
	}

	err = tx.SelectContext(ctx, &products, query, args...)
	return products, err
}
