package patagonia

import (
	"context"
	"net/http"

	"github.com/shubhamgptln/hotels-data-merge/app/domain/model"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/httpapi"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/httputil"
)

func (c *PatagoniaClient) FetchHotelsData(ctx context.Context) ([]*model.Hotel, error) {
	api := &getPatagoniaDataAPI{
		endpoint: c.getEndpoint,
	}
	receivedData := make([]*model.Hotel, 0)
	successful, err := c.httpClient.Call(context.Background(), api)
	if successful {
		for _, hotel := range api.resp {
			receivedData = append(receivedData, hotel.ClientHotelDataToDomainModel())
		}
	}
	return receivedData, err
}

var _ httpapi.API = new(getPatagoniaDataAPI)

type getPatagoniaDataAPI struct {
	endpoint string
	resp     []*HotelData
}

func (api *getPatagoniaDataAPI) BuildRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api.endpoint, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (api *getPatagoniaDataAPI) ParseResponse(_ context.Context, _ *http.Request, resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return httputil.BuildErrorWithResponse(httpapi.ErrResponseNotOk, resp)
	}
	if err := httputil.UnmarshalResponseBody(resp, &api.resp); err != nil {
		return err
	}
	return nil
}
