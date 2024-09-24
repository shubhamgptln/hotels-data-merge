package acme

import (
	"strconv"
	"strings"

	"hotel-data-merge/app/domain/model"
)

type HotelData struct {
	ID            string      `json:"Id"`
	DestinationID int64       `json:"DestinationId"`
	Name          string      `json:"Name"`
	Latitude      interface{} `json:"Latitude"`
	Longitude     interface{} `json:"Longitude"`
	Address       string      `json:"Address"`
	City          string      `json:"City"`
	Country       string      `json:"Country"`
	PostalCode    string      `json:"PostalCode"`
	Description   string      `json:"Description"`
	Facilities    []string    `json:"Facilities"`
}

func ClientHotelDataToDomainModel(data *HotelData) *model.Hotel {
	var lat, lng float64
	var ok bool
	//sanitize
	if data.Latitude != nil {
		switch v := data.Latitude.(type) {
		case string:
			lat, _ = strconv.ParseFloat(strings.TrimSpace(v), 64)
		default:
			lat, ok = v.(float64)
			if !ok {
				lat = 0
			}
		}
	}
	if data.Longitude != nil {
		switch v := data.Longitude.(type) {
		case string:
			lng, _ = strconv.ParseFloat(strings.TrimSpace(v), 64)
		default:
			lng, ok = v.(float64)
			if !ok {
				lat = 0
			}
		}
	}

	return &model.Hotel{
		ID:            data.ID,
		DestinationID: data.DestinationID,
		Name:          data.Name,
		Location: model.Location{
			Lat:     lat,
			Lng:     lng,
			Address: data.Address + "," + data.PostalCode,
			City:    data.City,
			Country: data.Country,
		},
		Description: data.Description,
		Amenities: model.Amenities{
			General: data.Facilities,
		},
	}
}
