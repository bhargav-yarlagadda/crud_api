package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

var movies = []Movie{
	{
		ID:    "1",
		Isbn:  "12345",
		Title: "Inception",
		Director: &Director{
			FirstName: "Christopher",
			LastName:  "Nolan",
		},
	},
}

func main() {
	movies = append(movies, Movie{
		ID:       "2",
		Isbn:     "2",
		Title:    "RRR",
		Director: &Director{FirstName: "Rajamouli", LastName: "SS"},
	})
	router := mux.NewRouter()
	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovieById).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovieById).Methods("DELETE")

	fmt.Println("starting server at port : 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Welcome to the movies server! Please use Postman to test the server.")
}

func deleteMovieById(w http.ResponseWriter, req *http.Request) {
	found := false
	params := mux.Vars(req)
	fmt.Println(params)
	var movie Movie
	for idx, item := range movies {
		if item.ID == params["id"] {
			found = true
			movie = item
			movies = append(movies[:idx], movies[idx+1:]...)
			break
		}
	}
	if !found {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	} else {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
			Movie   Movie  `json:"movie"`
		}{
			Message: "Movie deleted",
			Movie:   movie,
		})
	}
}

func updateMovieById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var updatedMovie Movie
	_ = json.NewDecoder(req.Body).Decode(&updatedMovie)

	found := false
	for idx, item := range movies {
		if item.ID == params["id"] {
			found = true
			// Update the movie details
			item.Title = updatedMovie.Title
			item.Isbn = updatedMovie.Isbn
			item.Director = updatedMovie.Director
			movies[idx] = item
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if found {
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
			Movie   Movie  `json:"movie"`
		}{
			Message: "Movie updated",
			Movie:   updatedMovie,
		})
	} else {
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	}
}

func createMovie(w http.ResponseWriter, req *http.Request) {
	var newMovie Movie
	_ = json.NewDecoder(req.Body).Decode(&newMovie)

	// Append the new movie to the movies slice
	movies = append(movies, newMovie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Message: "Movie created",
		Movie:   newMovie,
	})
}

func getMovieById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	found := false
	var movie Movie
	for _, item := range movies {
		if item.ID == params["id"] {
			found = true
			movie = item
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if found {
		json.NewEncoder(w).Encode(movie)
	} else {
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	}
}

func getMovies(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Fetching movies...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
