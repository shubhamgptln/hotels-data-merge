package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
	"github.com/shubhamgptln/hotels-data-merge/app/usecase"
	"github.com/shubhamgptln/hotels-data-merge/tests/mocks-gen/cachetest"
	"github.com/shubhamgptln/hotels-data-merge/tests/mocks-gen/extapitest"
	"github.com/shubhamgptln/hotels-data-merge/utils"
)

func TestGatherDataWithFiltering(t *testing.T) {
	type mocks struct {
		dataSupplierMock *extapitest.MockDataSupplier
		cacheMock        *cachetest.MockCaching
	}
	setup := func(t *testing.T) (deps mocks, procurer service.DataFetcher) {
		ctrl := gomock.NewController(t)
		deps.dataSupplierMock = extapitest.NewMockDataSupplier(ctrl)
		deps.cacheMock = cachetest.NewMockCaching(ctrl)
		procurer = usecase.NewDataProcurer(deps.cacheMock, deps.dataSupplierMock, deps.dataSupplierMock, deps.dataSupplierMock)
		return
	}
	stubErr := errors.New("unknown_error")
	hotelIds := []string{"iJhz", "f8c9"}
	destinationID := int64(5432)
	ctx := context.Background()
	mockData, _ := utils.GenerateFakeData(hotelIds, destinationID)
	hotels, _ := json.Marshal(mockData)

	t.Run("It should successfully procure data with cache hits", func(t *testing.T) {
		mocks, procurer := setup(t)
		mocks.cacheMock.EXPECT().Get(gomock.Any()).Return(string(hotels), true).Times(3)
		result, err := procurer.GatherDataWithFiltering(ctx, hotelIds, destinationID)
		require.NoError(t, err)
		for id, hotels := range result {
			require.NotNil(t, hotels)
			require.True(t, slices.Contains(hotelIds, id) || hotels[0].DestinationID == destinationID)
		}
	})

	t.Run("It should successfully procure data with cache miss", func(t *testing.T) {
		mocks, procurer := setup(t)
		mocks.cacheMock.EXPECT().Get(gomock.Any()).Return(gomock.Any(), false).Times(3)
		mocks.dataSupplierMock.EXPECT().FetchHotelsData(ctx).Return(mockData, nil).Times(3)
		mocks.cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), time.Hour).Return(true).Times(3)
		result, err := procurer.GatherDataWithFiltering(ctx, hotelIds, destinationID)
		require.NoError(t, err)
		for id, hotels := range result {
			require.NotNil(t, hotels)
			require.True(t, slices.Contains(hotelIds, id) || hotels[0].DestinationID == destinationID)
		}
	})

	t.Run("It should return error if supplier api returns error on cache miss", func(t *testing.T) {
		mocks, procurer := setup(t)
		mocks.cacheMock.EXPECT().Get(gomock.Any()).Return(gomock.Any(), false).Times(3)
		mocks.dataSupplierMock.EXPECT().FetchHotelsData(ctx).Return(nil, stubErr).Times(3)
		_, err := procurer.GatherDataWithFiltering(ctx, hotelIds, destinationID)
		require.EqualError(t, err, "ACME supplier returned error : unknown_error \nPatagonia supplier returned err : unknown_error \nPaperflies supplier returned err : unknown_error")
	})
}
