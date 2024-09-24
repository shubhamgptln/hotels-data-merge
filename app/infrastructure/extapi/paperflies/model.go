package paperflies

import "github.com/shubhamgptln/hotels-data-merge/app/domain/model"

type HotelData struct {
	ID                string        `json:"hotel_id"`
	DestinationID     int64         `json:"destination_id"`
	Name              string        `json:"hotel_name"`
	Location          Location      `json:"location"`
	Details           string        `json:"details"`
	Amenities         Amenities     `json:"amenities"`
	Images            PropertyImage `json:"images"`
	BookingConditions []string      `json:"booking_conditions"`
}

type Location struct {
	Address string `json:"address"`
	Country string `json:"country"`
}

type Amenities struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

type ImageDetails struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

type PropertyImage struct {
	Rooms []ImageDetails `json:"rooms"`
	Site  []ImageDetails `json:"site"`
}

func ClientHotelDataToDomainModel(data *HotelData) *model.Hotel {
	domainData := &model.Hotel{
		ID:            data.ID,
		DestinationID: data.DestinationID,
		Name:          data.Name,
		Location: model.Location{
			Address: data.Location.Address,
			Country: data.Location.Country,
		},
		Description: data.Details,
		Amenities: model.Amenities{
			General: data.Amenities.General,
			Room:    data.Amenities.Room,
		},
		Images: model.PropertyImages{
			Rooms: make([]model.ImageDetails, 0),
			Site:  make([]model.ImageDetails, 0),
		},
		BookingConditions: data.BookingConditions,
	}
	for _, roomImage := range data.Images.Rooms {
		domainData.Images.Rooms = append(domainData.Images.Rooms, ClientHotelImageToDomainImageModel(roomImage))
	}
	for _, siteImage := range data.Images.Site {
		domainData.Images.Site = append(domainData.Images.Site, ClientHotelImageToDomainImageModel(siteImage))
	}
	return domainData
}

func ClientHotelImageToDomainImageModel(data ImageDetails) model.ImageDetails {
	return model.ImageDetails{
		Link:        data.Link,
		Description: data.Caption,
	}
}
