package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler/view"
	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/rest"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
)

func FetchHotelsInformation(merger service.DataMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		hotelIDs := req.URL.Query()["hotel_ids"]
		destID := req.URL.Query().Get("destination_id")
		destinationID, err := strconv.ParseInt(destID, 10, 64)
		if hotelIDs == nil && (err != nil || destID == "") {
			rest.JSONResp(ctx, w, http.StatusBadRequest, &rest.JSONResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
				Debug:      fmt.Sprintf("error marshalling a req body: %v", req.Body),
			})
			return
		}

		hotels, err := merger.FetchBestHotelData(ctx, hotelIDs, destinationID)
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
