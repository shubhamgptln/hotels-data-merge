package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler/view"
	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/rest"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
)

func FetchHotelsInformation(merger service.DataMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var hotelFilterInput view.HotelDataRequest
		if err := json.NewDecoder(req.Body).Decode(&hotelFilterInput); err != nil {
			rest.JSONResp(ctx, w, http.StatusBadRequest, &rest.JSONResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
				Debug:      fmt.Sprintf("error marshalling a req body: %v", req.Body),
			})
			return
		}

		hotels, err := merger.FetchBestHotelData(ctx, hotelFilterInput.HotelIDs, hotelFilterInput.DestinationID)
		if err != nil {
			rest.JSONResp(ctx, w, http.StatusInternalServerError, &rest.JSONResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
				Debug:      fmt.Sprintf("error fetching supplier data: %v", hotels),
			})
			return
		}
		rest.JSONResp(ctx, w, http.StatusOK, &rest.JSONResponse{
			StatusCode: http.StatusOK,
			Data:       view.BuildHotelResponse(hotels),
		})
	}
}
