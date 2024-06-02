package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nyae44/cinerank/movies"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *movies.MovieAPIResponse
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage - 1
}
func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
func searchHandler(moviesapi *movies.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := moviesapi.FetchRatedMovies(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(moviesapi.PageSize))),
			Results:    results,
		}
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)

		fmt.Printf("%+v", results)
		fmt.Println("Search Query is: ", searchQuery)
		fmt.Println("Page is: ", page)
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")

	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"

	}
	apiKey := os.Getenv("RATED_MOVIES_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	moviesapi := movies.MovieClient(myClient, apiKey, 20)
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	mux.HandleFunc("/search", searchHandler(moviesapi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)

}
