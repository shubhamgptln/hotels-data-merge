package model

type Hotel struct {
	ID                string
	DestinationID     int64
	Name              string
	Location          Location
	Description       string
	Amenities         Amenities
	Images            PropertyImages
	BookingConditions []string
}

type Amenities struct {
	General []string
	Room    []string
}

type PropertyImages struct {
	Rooms     []ImageDetails
	Site      []ImageDetails
	Amenities []ImageDetails
}
