package service

import (
	"context"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
)

//go:generate mockgen -source $GOFILE -package servicetest -destination ../../../tests/mocks-gen/servicetest/procurer_mock.go

type DataFetcher interface {
	GatherDataWithFiltering(ctx context.Context, hotelIDs []string, destinationID int64) (map[string][]*model.Hotel, error)
}

type DataMerger interface {
	FetchBestHotelData(ctx context.Context, hotelIDs []string, destinationID int64) ([]*model.Hotel, error)
	MergeDataFromSuppliers(ctx context.Context, data map[string][]*model.Hotel) ([]*model.Hotel, error)
}
