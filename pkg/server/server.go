package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/trwndh/poc-online-store/pkg/logger"

	"github.com/go-chi/chi"
	orderHandler "github.com/trwndh/poc-online-store/module/order/handler"
	order "github.com/trwndh/poc-online-store/module/order/service"
	productHandler "github.com/trwndh/poc-online-store/module/product/handler"
	product "github.com/trwndh/poc-online-store/module/product/service"
)

type server struct {
	orderHandler   orderHandler.Handler
	productHandler productHandler.Handler
	router         *chi.Mux
	listenErrorCh  chan error
}

func NewServer(
	ctx context.Context,
	order order.OrderService,
	product product.ProductService,
) *server {
	orderHandler := orderHandler.NewHandler(order)
	productHandler := productHandler.NewHandler(product)
	return &server{
		orderHandler:   orderHandler,
		productHandler: productHandler,
		router:         chi.NewRouter(),
	}
}

func (s *server) Run() (err error) {
	log := logger.Log

	s.Route()

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVICE_PORT"),
		Handler: s.router,
	}

	log.Info(fmt.Sprintf("starting http server at port %s", srv.Addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("error start http server " + err.Error())
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// handle err
		log.Fatal("error shutting down http server " + err.Error())
	}
	return err
}
