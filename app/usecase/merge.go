package usecase

import (
	"context"
	"fmt"
	"reflect"

	"hotel-data-merge/app/domain/model"
	"hotel-data-merge/app/domain/service"
	"hotel-data-merge/bootstrap"
)

var _ service.DataMerger = new(DataMerger)

type DataMerger struct {
	config   *bootstrap.Config
	procurer service.DataFetcher
}

func NewDataMerger(config *bootstrap.Config, procurer service.DataFetcher) *DataMerger {
	return &DataMerger{config: config, procurer: procurer}
}

func (dm *DataMerger) FetchBestHotelData(ctx context.Context, hotelIDs []string, destinationID int64) ([]*model.Hotel, error) {
	data, err := dm.procurer.GatherDataWithFiltering(ctx, hotelIDs, destinationID)
	if err != nil {
		return nil, err
	}
	return dm.MergeDataFromSuppliers(ctx, data)
}

func (dm *DataMerger) MergeDataFromSuppliers(ctx context.Context, hotelMap map[string][]*model.Hotel) ([]*model.Hotel, error) {
	var mergedHotelResp []*model.Hotel
	for id, supplierData := range hotelMap {
		result := model.Hotel{
			ID:            id,
			DestinationID: supplierData[0].DestinationID,
		}
		val := reflect.ValueOf(result)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		typeOf := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fmt.Printf("Field: %v \t Value: %v\n", typeOf.Field(i).Name, val.Field(i).Interface())
			switch typeOf.Field(i).Name {
			case "Name":
				candidateValues := make([]string, 0)
				for _, entry := range supplierData {
					candidateValues = append(candidateValues, entry.Name)
				}
				result.Name = dm.strategySwitcherString(candidateValues, dm.config.DataMergeStrategy["Name"])[0]
			case "Description":
				candidateValues := make([]string, 0)
				for _, entry := range supplierData {
					candidateValues = append(candidateValues, entry.Description)
				}
				result.Description = dm.strategySwitcherString(candidateValues, dm.config.DataMergeStrategy["Description"])[0]
			case "BookingConditions":
				candidateValues := make([]string, 0)
				for _, entry := range supplierData {
					candidateValues = append(candidateValues, entry.BookingConditions...)
				}
				result.BookingConditions = dm.strategySwitcherString(candidateValues, dm.config.DataMergeStrategy["BookingConditions"])
			case "Amenities":
				candidateGeneralAmenitiesValues := make([]string, 0)
				for _, entry := range supplierData {
					candidateGeneralAmenitiesValues = append(candidateGeneralAmenitiesValues, entry.Amenities.General...)
				}
				candidateRoomAmenitiesValues := make([]string, 0)
				for _, entry := range supplierData {
					candidateRoomAmenitiesValues = append(candidateRoomAmenitiesValues, entry.Amenities.Room...)
				}
				result.Amenities.General = dm.strategySwitcherString(candidateGeneralAmenitiesValues, dm.config.DataMergeStrategy["Amenities.General"])
				result.Amenities.Room = dm.strategySwitcherString(candidateRoomAmenitiesValues, dm.config.DataMergeStrategy["Amenities.Room"])
			case "Location":
				candidateLocationValues := make([]interface{}, 0)
				for _, entry := range supplierData {
					candidateLocationValues = append(candidateLocationValues, entry.Location)
				}
				value := dm.strategySwitcherInterface(candidateLocationValues, "", dm.config.DataMergeStrategy["Location"])
				location, _ := value.(model.Location)
				result.Location = location
			case "Images":
				//longest description for unique images
				candidateAmenitiesImageValues := make(map[string]string, 0)
				candidateSiteImageValues := make(map[string]string, 0)
				candidateRoomImageValues := make(map[string]string, 0)
				result.Images.Site = make([]model.ImageDetails, 0)
				result.Images.Rooms = make([]model.ImageDetails, 0)
				result.Images.Amenities = make([]model.ImageDetails, 0)
				for _, entry := range supplierData {
					for _, img := range entry.Images.Amenities {
						if len(candidateAmenitiesImageValues[img.Link]) < len(img.Description) {
							candidateAmenitiesImageValues[img.Link] = img.Description
						}
					}
				}
				for _, entry := range supplierData {
					for _, img := range entry.Images.Site {
						if len(candidateSiteImageValues[img.Link]) < len(img.Description) {
							candidateSiteImageValues[img.Link] = img.Description
						}
					}
				}
				for _, entry := range supplierData {
					for _, img := range entry.Images.Rooms {
						if len(candidateRoomImageValues[img.Link]) < len(img.Description) {
							candidateRoomImageValues[img.Link] = img.Description
						}
					}
				}
				for url, entry := range candidateAmenitiesImageValues {
					result.Images.Amenities = append(result.Images.Amenities, model.ImageDetails{
						Link:        url,
						Description: entry,
					})
				}
				for url, entry := range candidateSiteImageValues {
					result.Images.Site = append(result.Images.Site, model.ImageDetails{
						Link:        url,
						Description: entry,
					})
				}
				for url, entry := range candidateRoomImageValues {
					result.Images.Rooms = append(result.Images.Rooms, model.ImageDetails{
						Link:        url,
						Description: entry,
					})
				}
			}
		}
		mergedHotelResp = append(mergedHotelResp, &result)
	}

	return mergedHotelResp, nil
}

// An attempt to hide internal field details from usecase layer implementation
//func (dm *DataMerger) collectAndMerge(supplierData []*model.Hotel, result *model.Hotel, strategy model.MergeStrategy) {
//	var fieldValueMap map[string][]interface{}
//	for _, entry := range supplierData {
//		val := reflect.ValueOf(entry).Elem()
//		typeOf := val.Type()
//		for i := 0; i < val.NumField(); i++ {
//			if val.Field(i).Kind() == reflect.String {
//				fieldValueMap[typeOf.Field(i).Name] = append(fieldValueMap[typeOf.Field(i).Name], val.Field(i).String())
//			} else if val.Field(i).Kind() == reflect.Float64 {
//				fieldValueMap[typeOf.Field(i).Name] = append(fieldValueMap[typeOf.Field(i).Name], val.Field(i).Float())
//			}
//		}
//	}
//	for _, data := range fieldValueMap {
//		result.Name = dm.strategySwitcherString(result, data, strategy)
//	}
//
//}

func (dm *DataMerger) strategySwitcherString(data []string, strategy model.MergeStrategy) []string {
	switch strategy {
	case model.Longest:
		return []string{strategy.LongestString(data)}
	case model.AppendUnique:
		return strategy.AppendUniqueEntries(data)
	}
	return nil
}

func (dm *DataMerger) strategySwitcherInterface(data []interface{}, fieldName string, strategy model.MergeStrategy) interface{} {
	switch strategy {
	case model.MajorityNonZero:
		return strategy.MajorityNonZeroField(data)
	case model.AppendUnique:
		return strategy.AssignFirstNonZeroField(data, fieldName)
	}
	return nil
}
