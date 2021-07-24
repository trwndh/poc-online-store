package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/trwndh/poc-online-store/module/order/model"
	"github.com/trwndh/poc-online-store/pkg/logger"
	"github.com/trwndh/poc-online-store/pkg/response"
)

func (h *handler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Log
	var noLocking bool
	locking := r.URL.Query().Get("no_locking")
	if locking == "true" {
		noLocking = true
	}
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
		Data model.Cart `json:"data"`
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

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "handler.checkout")
	defer span.Finish()

	regPaymentID := int64(0)
	if noLocking {
		regPaymentID, err = h.orderService.CheckoutNoLocking(ctx, request.Data)
	} else {
		regPaymentID, err = h.orderService.Checkout(ctx, request.Data)
	}
	if err != nil {
		response.Response(w, response.APIResponse{
			Data:  response.Empty{},
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	res := struct {
		RegisterPaymentID int64 `json:"register_payment_id"`
	}{
		RegisterPaymentID: regPaymentID,
	}
	response.Response(w, response.APIResponse{
		Data:  res,
		Code:  http.StatusCreated,
		Error: response.Empty{},
	})
	return
}
