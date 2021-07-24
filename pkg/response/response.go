package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  interface{} `json:"data"`
	Code  int         `json:"-"`
	Error interface{} `json:"error"`
}

type Empty struct{}
type Data struct{}

func Response(w http.ResponseWriter, res APIResponse) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	_ = json.NewEncoder(w).Encode(res)
}
