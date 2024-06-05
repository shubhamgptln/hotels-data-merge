package paperflies

import "github.com/shubhamgptln/hotels-data-merge/app/infrastructure/httpapi"

type PaperfliesClient struct {
	getEndpoint string
	httpClient  httpapi.Caller
}

func NewPaperfliesClient(
	getEndpoint string,
	httpClient httpapi.Caller,
) *PaperfliesClient {
	return &PaperfliesClient{
		getEndpoint: getEndpoint,
		httpClient:  httpClient,
	}
}
