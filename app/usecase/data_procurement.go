package usecase

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/exp/slices"

	"hotel-data-merge/app/domain/model"
	"hotel-data-merge/app/domain/service"
)

type DataProcurer struct {
	acmeClient       service.DataSupplier
	patagoniaClient  service.DataSupplier
	paperfliesClient service.DataSupplier
}

func NewDataProcurer(
	acmeClient service.DataSupplier,
	patagoniaClient service.DataSupplier,
	paperfliesClient service.DataSupplier) *DataProcurer {
	return &DataProcurer{
		acmeClient:       acmeClient,
		paperfliesClient: paperfliesClient,
		patagoniaClient:  patagoniaClient,
	}
}

func (dp *DataProcurer) GatherDataWithFiltering(ctx context.Context, hotelIDs []string, destinationID int64) (map[string][]*model.Hotel, error) {
	var wg sync.WaitGroup
	var acmeErr, paperfliesErr, patagoniaErr error
	var acmeData, paperfliesData, patagoniaData []*model.Hotel
	wg.Add(3)
	go func() error {
		defer wg.Done()
		acmeData, acmeErr = dp.acmeClient.FetchHotelsData(ctx)
		if acmeErr != nil {
			return acmeErr
		}
		acmeData = dp.filterDataBasedOnID(acmeData, hotelIDs, destinationID)
		return nil
	}()
	go func() error {
		defer wg.Done()
		paperfliesData, paperfliesErr = dp.paperfliesClient.FetchHotelsData(ctx)
		if paperfliesErr != nil {
			return paperfliesErr
		}
		paperfliesData = dp.filterDataBasedOnID(paperfliesData, hotelIDs, destinationID)
		return nil
	}()
	go func() error {
		defer wg.Done()
		patagoniaData, patagoniaErr = dp.patagoniaClient.FetchHotelsData(ctx)
		if patagoniaErr != nil {
			return patagoniaErr
		}
		patagoniaData = dp.filterDataBasedOnID(patagoniaData, hotelIDs, destinationID)
		return nil
	}()
	wg.Wait()
	fmt.Printf("length %v,%v,%v", len(patagoniaData), len(paperfliesData), len(acmeData))
	if acmeErr != nil && patagoniaErr != nil && paperfliesErr != nil {
		return nil, fmt.Errorf("ACME supplier returned error : %v \n"+
			"Patagonia supplier returned err : %v \n"+
			"Paperflies supplier returned err : %v", acmeErr, patagoniaErr, paperfliesErr)
	}
	var result map[string][]*model.Hotel
	for _, entry := range patagoniaData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	for _, entry := range acmeData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	for _, entry := range paperfliesData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	fmt.Println(len(result))
	return result, nil
}

func (dp *DataProcurer) filterDataBasedOnID(data []*model.Hotel, hotelIDs []string, destinationID int64) []*model.Hotel {
	result := make([]*model.Hotel, 0)
	if len(hotelIDs) > 0 {
		for _, hotel := range data {
			if slices.Contains(hotelIDs, hotel.ID) {
				result = append(result, hotel)
			}
		}
		return result
	}
	for _, hotel := range data {
		if hotel.DestinationID == destinationID {
			result = append(result, hotel)
		}
	}
	return result
}
