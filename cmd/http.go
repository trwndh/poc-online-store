package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	productMySQLRepo "github.com/trwndh/poc-online-store/module/product/repository/mysql"
	product "github.com/trwndh/poc-online-store/module/product/service"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/trwndh/poc-online-store/config"
	orderMySQLRepo "github.com/trwndh/poc-online-store/module/order/repository/mysql"
	order "github.com/trwndh/poc-online-store/module/order/service"
	"github.com/trwndh/poc-online-store/pkg/logger"
	"github.com/trwndh/poc-online-store/pkg/server"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var StartHTTP = &cobra.Command{
	Use:   "http-start",
	Short: "Start http server for REST API",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.Log

		err := config.Setup("config.env")
		if err != nil {
			log.Fatal("error when loading config", zap.Error(err))
			return
		}

		jcfg, _ := jaegercfg.FromEnv()
		tracer, closer, err := jcfg.NewTracer()
		if err != nil {
			log.Fatal("failed to init tracer", zap.Error(err))
		}
		defer func() { _ = closer.Close() }()
		opentracing.SetGlobalTracer(tracer)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB"))
		db, err := sqlx.Connect(os.Getenv("DB_DRIVER"), dsn)
		if err != nil {
			log.Fatal("failed to connect database", zap.Error(err))
			time.Sleep(10 * time.Second)
			for err != nil {
				db, err = sqlx.Connect(os.Getenv("DB_DRIVER"), dsn)
				if err == nil {
					break
				}
			}
		}

		db.SetConnMaxIdleTime(0)
		db.SetConnMaxLifetime(0)
		db.SetMaxOpenConns(0)
		log.Info(fmt.Sprintf("connected to %s database at %s:%s", db.DriverName(), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")))

		ctx := context.Background()

		orderRepo := orderMySQLRepo.NewOrderRepository(ctx, db)
		orderService := order.NewOrderService(ctx, orderRepo)

		productRepo := productMySQLRepo.NewProductRepository(ctx, db)
		productService := product.NewProductService(ctx, productRepo)
		srv := server.NewServer(ctx, orderService, productService)
		srv.Run()
	},
}
