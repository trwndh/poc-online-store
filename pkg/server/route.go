package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (s *server) Route() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)
	s.router.Use(Trace)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			response.Response(w, response.APIResponse{
				Data:  "API V1",
				Code:  http.StatusOK,
				Error: response.Empty{},
			})
		})
		r.Post("/order/checkout", s.orderHandler.CheckoutHandler)

		r.Route("/product", func(r chi.Router) {
			r.Get("/", s.productHandler.GetProducts)
			r.Get("/{id}", s.productHandler.GetProductByID)
			r.Put("/stock", s.productHandler.UpdateProductStock)
			r.Post("/", s.productHandler.AddProduct)
		})
	})
}
