package movies

import "net/http"

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

func MovieClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
