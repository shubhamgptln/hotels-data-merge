package patagonia

import "hotel-data-merge/app/domain/model"

type HotelData struct {
	ID            string        `json:"id"`
	DestinationID int64         `json:"destination"`
	Name          string        `json:"name"`
	Lat           float64       `json:"lat"`
	Lng           float64       `json:"lng"`
	Address       string        `json:"address"`
	Information   string        `json:"info"`
	Amenities     []string      `json:"amenities"`
	Images        PropertyImage `json:"images"`
}

type ImageDetails struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type PropertyImage struct {
	Rooms     []ImageDetails `json:"rooms"`
	Amenities []ImageDetails `json:"amenities"`
}

func ClientHotelDataToDomainModel(data *HotelData) *model.Hotel {
	domainData := &model.Hotel{
		ID:            data.ID,
		DestinationID: data.DestinationID,
		Name:          data.Name,
		Location: model.Location{
			Address: data.Address,
			Lat:     data.Lat,
			Lng:     data.Lng,
		},
		Description: data.Information,
		Amenities: model.Amenities{
			General: data.Amenities,
		},
		Images: model.PropertyImages{
			Rooms:     make([]model.ImageDetails, 0),
			Amenities: make([]model.ImageDetails, 0),
		},
	}
	for _, roomImage := range data.Images.Rooms {
		domainData.Images.Rooms = append(domainData.Images.Rooms, ClientHotelImageToDomainImageModel(roomImage))
	}
	for _, siteImage := range data.Images.Amenities {
		domainData.Images.Amenities = append(domainData.Images.Amenities, ClientHotelImageToDomainImageModel(siteImage))
	}
	return domainData
}

func ClientHotelImageToDomainImageModel(data ImageDetails) model.ImageDetails {
	return model.ImageDetails{
		Link:        data.Url,
		Description: data.Description,
	}
}
