package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Debug      interface{} `json:"debug,omitempty"`
}

func JSONResp(ctx context.Context, w http.ResponseWriter, statusCode int, body *JSONResponse) {
	// NOTE: Must call w.Header().Set before calling w.WriteHeader, else Content-Type will be text/plain
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	b, err := json.Marshal(body)
	if err != nil {
		JSONResp(ctx, w, http.StatusInternalServerError, &JSONResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
			Debug:      fmt.Sprintf("context: marshalling a json body: %v", body),
		})
		return
	}

	count, err := w.Write(b)
	if count == 0 || err != nil {
		JSONResp(ctx, w, http.StatusInternalServerError, &JSONResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
			Debug:      fmt.Sprintf("context: writing a json body: %v with bytes_count=%d", body, count),
		})
		return
	}
}
