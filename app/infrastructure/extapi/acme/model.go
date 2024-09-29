package acme

import (
	"strconv"
	"strings"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
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

func (hd *HotelData) ClientHotelDataToDomainModel() *model.Hotel {
	var lat, lng float64
	var ok bool
	//sanitize
	if hd.Latitude != nil {
		switch v := hd.Latitude.(type) {
		case string:
			lat, _ = strconv.ParseFloat(strings.TrimSpace(v), 64)
		default:
			lat, ok = v.(float64)
			if !ok {
				lat = 0
			}
		}
	}
	//sanitize
	if hd.Longitude != nil {
		switch v := hd.Longitude.(type) {
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
		ID:            hd.ID,
		DestinationID: hd.DestinationID,
		Name:          hd.Name,
		Location: model.Location{
			Lat:     lat,
			Lng:     lng,
			Address: hd.Address + "," + hd.PostalCode,
			City:    hd.City,
			Country: hd.Country,
		},
		Description: hd.Description,
		Amenities: model.Amenities{
			General: hd.Facilities,
		},
	}
}
