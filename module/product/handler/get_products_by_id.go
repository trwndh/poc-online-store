package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/opentracing/opentracing-go"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (h *handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "handler.GetProducts")
	defer span.Finish()
	ID := chi.URLParam(r, "id")
	if ID == "" {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusNotAcceptable,
			Error: "required parameter: id",
		})
		return
	}

	productID, err := strconv.Atoi(ID)
	if err != nil {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	products, err := h.productService.GetProductByID(ctx, int64(productID))
	if err != nil {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	response.Response(w, response.APIResponse{
		Data:  products,
		Code:  http.StatusOK,
		Error: response.Empty{},
	})
	return
}
