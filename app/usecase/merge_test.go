package usecase_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/service"
	"github.com/shubhamgptln/hotels-data-merge/app/usecase"
	"github.com/shubhamgptln/hotels-data-merge/bootstrap"
	"github.com/shubhamgptln/hotels-data-merge/tests/mocks-gen/servicetest"
	"github.com/shubhamgptln/hotels-data-merge/utils"
)

func TestFetchBestHotelData(t *testing.T) {
	type mocks struct {
		dataFetcherMock *servicetest.MockDataFetcher
		dataMergerMock  *servicetest.MockDataMerger
	}
	setup := func(t *testing.T) (deps mocks, merger service.DataMerger) {
		ctrl := gomock.NewController(t)
		deps.dataFetcherMock = servicetest.NewMockDataFetcher(ctrl)
		deps.dataMergerMock = servicetest.NewMockDataMerger(ctrl)
		merger = usecase.NewDataMerger(&bootstrap.Config{}, deps.dataFetcherMock)
		return
	}
	stubErr := errors.New("unknown_error")
	hotelIds := []string{"iJhz", "f8c9"}
	destinationID := int64(5432)
	ctx := context.Background()
	mockData, _ := utils.GenerateFakeData(hotelIds, destinationID)
	mockData2, _ := utils.GenerateFakeData(hotelIds, destinationID)
	mockData3, _ := utils.GenerateFakeData(hotelIds, destinationID)

	mockResult := map[string][]*model.Hotel{
		mockData[0].ID: {mockData[0], mockData2[0], mockData3[0]},
		mockData[1].ID: {mockData[1], mockData2[1], mockData3[1]},
		mockData[2].ID: {mockData[2], mockData2[2], mockData3[2]},
	}

	t.Run("It should successfully procure and merge data", func(t *testing.T) {
		mocks, merger := setup(t)
		mocks.dataFetcherMock.EXPECT().GatherDataWithFiltering(gomock.Any(), hotelIds, destinationID).Return(mockResult, nil)
		result, err := merger.FetchBestHotelData(ctx, hotelIds, destinationID)
		require.NoError(t, err)
		require.NotNil(t, result)
		for _, hotel := range result {
			require.True(t, slices.Contains(hotelIds, hotel.ID) || hotel.DestinationID == destinationID)
		}
	})

	t.Run("It should return error if fetcher returns error", func(t *testing.T) {
		mocks, merger := setup(t)
		mocks.dataFetcherMock.EXPECT().GatherDataWithFiltering(gomock.Any(), hotelIds, destinationID).Return(nil, stubErr)
		result, err := merger.FetchBestHotelData(ctx, hotelIds, destinationID)
		require.EqualError(t, err, stubErr.Error())
		require.Nil(t, result)
	})
}
