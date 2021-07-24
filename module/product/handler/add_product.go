package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	"github.com/trwndh/poc-online-store/module/product/model"
	"github.com/trwndh/poc-online-store/pkg/logger"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (h *handler) AddProduct(w http.ResponseWriter, r *http.Request) {
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
		Data model.Product `json:"data"`
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

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "handler.AddProduct")
	defer span.Finish()
	logger.Log.Info("request", zap.Any("value", request.Data))
	product, err := h.productService.AddProduct(ctx, request.Data)
	if err != nil {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	response.Response(w, response.APIResponse{
		Data:  product,
		Code:  http.StatusCreated,
		Error: response.Empty{},
	})
	return
}
