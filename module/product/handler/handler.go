package handler

import (
	"net/http"

	"github.com/trwndh/poc-online-store/module/product/service"
)

type handler struct {
	productService service.ProductService
}

type Handler interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProductByID(w http.ResponseWriter, r *http.Request)
	UpdateProductStock(w http.ResponseWriter, r *http.Request)
	AddProduct(w http.ResponseWriter, r *http.Request)
}

func NewHandler(productService service.ProductService) Handler {
	return &handler{productService: productService}
}
