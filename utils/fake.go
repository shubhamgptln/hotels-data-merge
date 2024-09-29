package utils

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler/view"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
)

func GenerateFakeData(hotelIDs []string, destinationID int64) ([]*model.Hotel, string) {
	hotels := make([]*model.Hotel, 0)
	for _, id := range hotelIDs {
		var tempHotel model.Hotel
		_ = gofakeit.Struct(&tempHotel)
		tempHotel.ID = id
		hotels = append(hotels, &tempHotel)
	}
	var tempDestHotel model.Hotel
	_ = gofakeit.Struct(&tempDestHotel)
	tempDestHotel.DestinationID = destinationID
	hotels = append(hotels, &tempDestHotel)
	resp := `{
		"status_code":200,
		"data": %s
	}`
	bytes, _ := json.Marshal(view.BuildHotelResponse(hotels))
	return hotels, fmt.Sprintf(resp, string(bytes))
}
