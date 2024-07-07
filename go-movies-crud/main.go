package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MIGUELNINOSILVA/go-projects/go-movies-crud/middlewares"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieChan := make(chan *Movie)

	go func() {
		for _, item := range movies {
			if item.ID == params["id"] {
				movieChan <- &item
				return
			}
		}
		movieChan <- nil
	}()

	movie := <-movieChan
	if movie != nil {
		json.NewEncoder(w).Encode(movie)
	} else {
		http.Error(w, "Movie not found", http.StatusNotFound)
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	var movie Movie
	for index, item := range movies {
		if item.ID == params["id"] {
			movie = movies[index]
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	fmt.Println(movie)

	json.NewEncoder(w).Encode(movie)
}

func createtMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		fmt.Println(err)
		return
	}
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func main() {
	r := mux.NewRouter()
	r.Use(middlewares.SetHttpHeaders)

	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "438227",
		Title: "Movie one",
		Director: &Director{
			FirstName: "John",
			LastName:  "Doe",
		},
	},
	)

	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "438227",
		Title: "Movie two",
		Director: &Director{
			FirstName: "Steve",
			LastName:  "Smith",
		},
	},
	)

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createtMovie).Methods("POST")
	// r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8081\n")
	log.Fatal(http.ListenAndServe(":8081", r))
}
