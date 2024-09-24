package service

import (
	"context"

	"hotel-data-merge/app/domain/model"
)

type DataFetcher interface {
	GatherDataWithFiltering(ctx context.Context, hotelIDs []string, destinationID int64) (map[string][]*model.Hotel, error)
}

type DataMerger interface {
	FetchBestHotelData(ctx context.Context, hotelIDs []string, destinationID int64) ([]*model.Hotel, error)
	MergeDataFromSuppliers(ctx context.Context, data map[string][]*model.Hotel) ([]*model.Hotel, error)
}
