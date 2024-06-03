package service

import (
	"context"

	"hotel-data-merge/app/domain/model"
)

type DataSupplier interface {
	FetchHotelsData(ctx context.Context) ([]*model.Hotel, error)
}
