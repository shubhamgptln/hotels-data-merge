package acme

import "hotel-data-merge/app/infrastructure/httpapi"

type ACMEClient struct {
	getEndpoint string
	httpClient  httpapi.Caller
}

func NewACMEClient(
	getEndpoint string,
	httpClient httpapi.Caller,
) *ACMEClient {
	return &ACMEClient{
		getEndpoint: getEndpoint,
		httpClient:  httpClient,
	}
}
