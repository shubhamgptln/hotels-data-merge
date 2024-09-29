package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/cache"
)

const (
	patagoniaCacheKey  = "patagonia"
	paperfliesCacheKey = "paperflies"
	acmeCacheKey       = "acme"
)

type DataProcurer struct {
	cache            cache.Caching
	acmeClient       service.DataSupplier
	patagoniaClient  service.DataSupplier
	paperfliesClient service.DataSupplier
}

func NewDataProcurer(
	inMemoryCache cache.Caching,
	acmeClient service.DataSupplier,
	patagoniaClient service.DataSupplier,
	paperfliesClient service.DataSupplier) *DataProcurer {
	return &DataProcurer{
		cache:            inMemoryCache,
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
		value, exists := dp.cache.Get(acmeCacheKey)
		if !exists {
			acmeData, acmeErr = dp.cacheMissFunc(ctx, dp.acmeClient, time.Hour)
		} else {
			acmeData, acmeErr = dp.cacheHitFunc(ctx, value)
		}
		if acmeErr != nil {
			return acmeErr
		}
		acmeData = dp.filterDataBasedOnID(acmeData, hotelIDs, destinationID)
		return nil
	}()
	go func() error {
		defer wg.Done()
		value, exists := dp.cache.Get(paperfliesCacheKey)
		if !exists {
			paperfliesData, paperfliesErr = dp.cacheMissFunc(ctx, dp.paperfliesClient, time.Hour)
		} else {
			paperfliesData, paperfliesErr = dp.cacheHitFunc(ctx, value)
		}
		if paperfliesErr != nil {
			return paperfliesErr
		}
		paperfliesData = dp.filterDataBasedOnID(paperfliesData, hotelIDs, destinationID)
		return nil
	}()
	go func() error {
		defer wg.Done()
		value, exists := dp.cache.Get(patagoniaCacheKey)
		if !exists {
			patagoniaData, patagoniaErr = dp.cacheMissFunc(ctx, dp.patagoniaClient, time.Hour)
		} else {
			patagoniaData, patagoniaErr = dp.cacheHitFunc(ctx, value)
		}
		if patagoniaErr != nil {
			return patagoniaErr
		}
		patagoniaData = dp.filterDataBasedOnID(patagoniaData, hotelIDs, destinationID)
		return nil
	}()
	wg.Wait()
	if acmeErr != nil && patagoniaErr != nil && paperfliesErr != nil {
		return nil, fmt.Errorf("ACME supplier returned error : %v \n"+
			"Patagonia supplier returned err : %v \n"+
			"Paperflies supplier returned err : %v", acmeErr, patagoniaErr, paperfliesErr)
	}
	result := make(map[string][]*model.Hotel, 0)
	for _, entry := range patagoniaData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	for _, entry := range acmeData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	for _, entry := range paperfliesData {
		result[entry.ID] = append(result[entry.ID], entry)
	}
	return result, nil
}

func (dp *DataProcurer) filterDataBasedOnID(data []*model.Hotel, hotelIDs []string, destinationID int64) []*model.Hotel {
	result := make([]*model.Hotel, 0)
	for _, hotel := range data {
		if hotel.DestinationID == destinationID {
			result = append(result, hotel)
		}
	}
	if len(hotelIDs) > 0 {
		for _, hotel := range data {
			if slices.Contains(hotelIDs, hotel.ID) && hotel.DestinationID != destinationID {
				result = append(result, hotel)
			}
		}
	}
	return result
}

func (dp *DataProcurer) cacheMissFunc(ctx context.Context, client service.DataSupplier, expiry time.Duration) (
	[]*model.Hotel, error) {
	data, err := client.FetchHotelsData(ctx)
	if err != nil {
		return nil, err
	}
	byteData, err := json.Marshal(data)
	if err == nil {
		dp.cache.Set(acmeCacheKey, string(byteData), expiry)
	}
	return data, nil
}

func (dp *DataProcurer) cacheHitFunc(ctx context.Context, value interface{}) (
	[]*model.Hotel, error) {
	var data []*model.Hotel
	val, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("value of unidentified type")
	}
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
