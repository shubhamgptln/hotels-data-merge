package service

import (
	"context"

	"hotel-data-merge/app/domain/model"
)

type DataFetcher interface {
	GatherDataWithFiltering(ctx context.Context, hotelIDs []string, destinationID int64) ([]*model.Hotel, error)
}
