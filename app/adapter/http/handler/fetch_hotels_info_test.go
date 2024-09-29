package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler"
	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler/view"
	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
	"github.com/shubhamgptln/hotels-data-merge/tests/mocks-gen/servicetest"
)

func TestFetchHotelData(t *testing.T) {
	type mocks struct {
		dataMergerMock *servicetest.MockDataMerger
	}
	hotelId := []string{"iJhz", "f8c9"}
	destination := "5432"
	setup := func(t *testing.T, getURL string, hotelIDs []string, destinationID string, withoutQueryParam bool) (mocks mocks, s *httptest.Server, url string, closeFn func()) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		r := chi.NewRouter()
		mocks.dataMergerMock = servicetest.NewMockDataMerger(ctrl)
		r.Use()
		r.Get("/hotels-data/filter", handler.FetchHotelsInformation(mocks.dataMergerMock))
		s = httptest.NewServer(r)
		if !withoutQueryParam {
			url = fmt.Sprintf(getURL, s.URL, strings.Join(hotelIDs[:], ","), destinationID)
		} else {
			url = fmt.Sprintf(getURL, s.URL)
		}
		closeFn = func() { ctrl.Finish(); s.Close() }
		return
	}

	hotels, respJSON := generateFakeData(hotelId, destination)
	stubError := errors.New("stub error")
	urlWithParam := "%s/hotels-data/filter?hotel_ids=%v&destination_id=%s"
	urlWithoutParam := "%s/hotels-data/filter"

	t.Run("success", func(t *testing.T) {
		mocks, s, url, closeFn := setup(t, urlWithParam, hotelId, destination, false)
		defer closeFn()
		mocks.dataMergerMock.EXPECT().
			FetchBestHotelData(gomock.Any(), hotelId, int64(5432)).
			Times(1).
			Return(hotels, nil)
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := s.Client().Do(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var b bytes.Buffer
		_, _ = b.ReadFrom(resp.Body)

		require.JSONEq(
			t,
			respJSON,
			b.String(),
		)

	})

	t.Run("fetcher function return error", func(t *testing.T) {
		mocks, s, url, closeFn := setup(t, urlWithParam, hotelId, destination, false)
		defer closeFn()
		mocks.dataMergerMock.EXPECT().
			FetchBestHotelData(gomock.Any(), hotelId, int64(5432)).
			Times(1).
			Return(nil, stubError)

		req, _ := http.NewRequest("GET", url, nil)
		resp, err := s.Client().Do(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var b bytes.Buffer
		_, _ = b.ReadFrom(resp.Body)
		assert.JSONEq(
			t,
			`{ "status_code": 500, "error": "stub error", "debug": "error from fetcher"}`,
			b.String(),
		)
	})

	t.Run("hotelID and destination_id both empty", func(t *testing.T) {
		_, s, url, closeFn := setup(t, urlWithoutParam, nil, "", true)
		defer closeFn()

		req, _ := http.NewRequest("GET", url, nil)
		resp, err := s.Client().Do(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b bytes.Buffer
		_, _ = b.ReadFrom(resp.Body)
		assert.JSONEq(
			t,
			`{ "status_code": 400, "error": "BAD_REQUEST" , "debug":"error getting req params"}`,
			b.String(),
		)
	})
}

func generateFakeData(hotelIDs []string, destinationID string) ([]*model.Hotel, string) {
	hotels := make([]*model.Hotel, 0)
	for _, id := range hotelIDs {
		var tempHotel model.Hotel
		_ = gofakeit.Struct(&tempHotel)
		tempHotel.ID = id
		hotels = append(hotels, &tempHotel)
	}
	var tempDestHotel model.Hotel
	_ = gofakeit.Struct(&tempDestHotel)
	tempDestHotel.ID = destinationID
	hotels = append(hotels, &tempDestHotel)
	resp := `{
		"status_code":200,
		"data": %s
	}`
	bytes, _ := json.Marshal(view.BuildHotelResponse(hotels))
	return hotels, fmt.Sprintf(resp, string(bytes))
}
