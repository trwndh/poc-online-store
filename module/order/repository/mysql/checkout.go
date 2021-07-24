package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/trwndh/poc-online-store/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/order/model"
)

func (o *OrderRepo) Checkout(ctx context.Context, userID int64, cartID int64, cartItems []model.CartItem) (registerPaymentID int64, err error) {
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
	products, err := findProduct(ctx, tx, productIDs)
	if err != nil {
		return 0, err
	}
	if len(products) == 0 {
		return 0, errors.New("error product not found:" + fmt.Sprintf(" %+v", cartItems))
	}

	// 2. if quantity avail >= quantity requested, continue, else, return error with stock unavailable.

	for _, product := range products {
		if qtyRequested, ok := mapProductStock[product.ID]; ok {
			logger.Info("updating product", zap.Any("product", product.Name))
			if qtyRequested <= product.Stock {
				err = updateProductStock(ctx, tx, product.ID, qtyRequested)
				if err != nil {
					return 0, err
				}
			} else {
				errMsg := fmt.Sprintf("insufficient stock for product %s, want %d but only available %d", product.Name, qtyRequested, product.Stock)
				return 0, errors.New(errMsg)
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

func findProduct(ctx context.Context, tx *sqlx.Tx, productIDs []int64) (products []model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.order.checkout.findProduct")
	defer span.Finish()

	query, args, err := sqlx.In("SELECT id, name, stock FROM product WHERE id IN (?) FOR UPDATE", productIDs)
	if err != nil {
		return []model.Product{}, err
	}

	err = tx.SelectContext(ctx, &products, query, args...)
	return products, err
}

func updateProductStock(ctx context.Context, tx *sqlx.Tx, productID int64, qty int32) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.order.checkout.updateProductStock")
	defer span.Finish()

	query := "UPDATE product SET stock = stock - ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, qty, productID)

	return err
}

func registerPayment(ctx context.Context, tx *sqlx.Tx, cartID int64, userID int64, jsonCartItem string) (registerPaymentID int64, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.order.checkout.registerPayment")
	defer span.Finish()

	query := "INSERT INTO register_payment(user_id, cart_id, product, created_at, expired_at) VALUES (?,?,?,?,?)"
	now := time.Now().Local()
	expireAt := now.Add(3 * time.Minute)

	res, err := tx.ExecContext(ctx, query, userID, cartID, jsonCartItem, now.Format("2006-01-02 15:04:05"), expireAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertID, nil
}
