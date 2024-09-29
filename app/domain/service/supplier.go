package service

import (
	"context"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
)

//go:generate mockgen -source $GOFILE -package extapitest -destination ../../../tests/mocks-gen/extapitest/supplier_mock.go

type DataSupplier interface {
	FetchHotelsData(ctx context.Context) ([]*model.Hotel, error)
}
