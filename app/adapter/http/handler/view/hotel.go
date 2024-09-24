package view

import "hotel-data-merge/app/domain/model"

type HotelDataRequest struct {
	DestinationID int64    `json:"destination_id"`
	HotelIDs      []string `json:"hotel_id"`
}

type Hotel struct {
	ID                string        `json:"id"`
	DestinationID     int64         `json:"destination_id"`
	Name              string        `json:"name"`
	Location          Location      `json:"location"`
	Description       string        `json:"description"`
	Amenities         Amenities     `json:"amenities"`
	Images            PropertyImage `json:"images"`
	BookingConditions []string      `json:"booking_conditions"`
}

type Location struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	Country string  `json:"country"`
}

type Amenities struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

type PropertyImage struct {
	Rooms     []ImageDetails `json:"rooms"`
	Site      []ImageDetails `json:"site"`
	Amenities []ImageDetails `json:"amenities"`
}

type ImageDetails struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}

func BuildHotelResponse(data []*model.Hotel) []*Hotel {
	var hotelView []*Hotel
	for _, entry := range data {
		hotelView = append(hotelView, &Hotel{
			ID:            entry.ID,
			DestinationID: entry.DestinationID,
			Name:          entry.Name,
			Location: Location{
				Lat:     entry.Location.Lat,
				Lng:     entry.Location.Lng,
				Address: entry.Location.Address,
				City:    entry.Location.City,
				Country: entry.Location.Country,
			},
			Description: entry.Description,
			Amenities: Amenities{
				General: entry.Amenities.General,
				Room:    entry.Amenities.Room,
			},
			Images: PropertyImage{
				Rooms:     ModelToViewImage(entry.Images.Rooms),
				Site:      ModelToViewImage(entry.Images.Site),
				Amenities: ModelToViewImage(entry.Images.Amenities),
			},
			BookingConditions: entry.BookingConditions,
		})
	}
	return hotelView
}

func ModelToViewImage(data []model.ImageDetails) []ImageDetails {
	var result []ImageDetails
	for _, image := range data {
		result = append(result, ImageDetails{
			Link:        image.Link,
			Description: image.Description,
		})
	}
	return result
}
