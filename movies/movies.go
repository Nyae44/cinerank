package movies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Movie struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	ID               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_avarage"`
	VoteCount        int     `json:"vote_count"`
}

type MovieAPIResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

func (c *Client) FetchRatedMovies(query, page string) (*MovieAPIResponse, error) {
	endpoint := fmt.Sprintf("https://api.themoviedb.org/3/movie/top_rated?api_key=%s&language=en-US&page=%s", url.QueryEscape(c.key), url.QueryEscape(page))

	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	res := &MovieAPIResponse{}
	return res, json.Unmarshal(body, res)
}

func MovieClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
