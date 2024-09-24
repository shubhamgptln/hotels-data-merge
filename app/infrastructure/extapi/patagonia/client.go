package patagonia

import "hotel-data-merge/app/infrastructure/httpapi"

type PatagoniaClient struct {
	getEndpoint string
	httpClient  httpapi.Caller
}

func NewPatagoniaClient(
	getEndpoint string,
	httpClient httpapi.Caller,
) *PatagoniaClient {
	return &PatagoniaClient{
		getEndpoint: getEndpoint,
		httpClient:  httpClient,
	}
}
