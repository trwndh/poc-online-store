package handler

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "handler.GetProducts")
	defer span.Finish()

	products, err := h.productService.GetProducts(ctx)
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
