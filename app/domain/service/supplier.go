package service

import (
	"context"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
)

type DataSupplier interface {
	FetchHotelsData(ctx context.Context) ([]*model.Hotel, error)
}
