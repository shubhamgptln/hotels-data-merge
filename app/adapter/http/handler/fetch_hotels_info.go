package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler/view"
	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/rest"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
)

func FetchHotelsInformation(merger service.DataMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		hotelIDs := req.URL.Query()["hotel_ids"]
		if hotelIDs != nil {
			hotelIDs = strings.Split(hotelIDs[0], ",")
		}
		destID := req.URL.Query().Get("destination_id")
		destinationID, err := strconv.ParseInt(destID, 10, 64)
		if hotelIDs == nil && (err != nil || destID == "") {
			rest.JSONResp(ctx, w, http.StatusBadRequest, &rest.JSONResponse{
				StatusCode: http.StatusBadRequest,
				Error:      fmt.Errorf("BAD_REQUEST").Error(),
				Debug:      fmt.Sprintf("error getting req params"),
			})
			return
		}

		hotels, err := merger.FetchBestHotelData(ctx, hotelIDs, destinationID)
		if err != nil {
			rest.JSONResp(ctx, w, http.StatusInternalServerError, &rest.JSONResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
				Debug:      fmt.Sprintf("error from fetcher"),
			})
			return
		}
		rest.JSONResp(ctx, w, http.StatusOK, &rest.JSONResponse{
			StatusCode: http.StatusOK,
			Data:       view.BuildHotelResponse(hotels),
		})
	}
}
