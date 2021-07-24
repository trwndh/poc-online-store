package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/trwndh/poc-online-store/pkg/logger"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (h *handler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	log := logger.Log

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("fail to read body " + err.Error())
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusNotAcceptable,
			Error: "Failed to read body",
		})
		return
	}

	request := struct {
		Data struct {
			ProductID int64 `json:"product_id"`
			Stock     int32 `json:"stock"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Error(err.Error())
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "handler.UpdateProductStock")
	defer span.Finish()

	err = h.productService.UpdateProductStock(ctx, request.Data.ProductID, request.Data.Stock)
	if err != nil {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	response.Response(w, response.APIResponse{
		Data:  "Successfully update product stock",
		Code:  http.StatusCreated,
		Error: response.Empty{},
	})
	return
}
