package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// Director struct represents the director of a movie
type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Movie struct represents the details of a movie
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"` // Director is a pointer to the Director struct
}

// movies slice holds the list of movies in memory
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
	// Adding another movie to the initial movies slice
	movies = append(movies, Movie{
		ID:       "2",
		Isbn:     "2",
		Title:    "RRR",
		Director: &Director{FirstName: "Rajamouli", LastName: "SS"},
	})

	// Create a new router using the mux package
	router := mux.NewRouter()

	// Define route handlers
	router.HandleFunc("/", welcome).Methods("GET")             // Welcome endpoint
	router.HandleFunc("/movies", getMovies).Methods("GET")     // Get all movies
	router.HandleFunc("/movies/{id}", getMovieById).Methods("GET") // Get movie by ID
	router.HandleFunc("/movies", createMovie).Methods("POST")  // Create a new movie
	router.HandleFunc("/movies/{id}", updateMovieById).Methods("PUT") // Update a movie by ID
	router.HandleFunc("/movies/{id}", deleteMovieById).Methods("DELETE") // Delete a movie by ID

	// Start the server on port 8080
	fmt.Println("starting server at port : 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// welcome function handles the GET request to the root endpoint
func welcome(w http.ResponseWriter, req *http.Request) {
	// Sending a welcome message as a response
	fmt.Fprint(w, "Welcome to the movies server! Please use Postman to test the server.")
}

// deleteMovieById function handles the DELETE request to delete a movie by ID
func deleteMovieById(w http.ResponseWriter, req *http.Request) {
	// Check if movie ID exists
	found := false
	params := mux.Vars(req) // Get the movie ID from URL
	var movie Movie
	for idx, item := range movies {
		if item.ID == params["id"] {
			// Movie found, deleting it from the slice
			found = true
			movie = item
			movies = append(movies[:idx], movies[idx+1:]...) // Remove the movie by index
			break
		}
	}

	// Return appropriate response
	w.Header().Set("Content-Type", "application/json")
	if !found {
		// If movie not found, return error message
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	} else {
		// If movie is deleted, return success message
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
			Movie   Movie  `json:"movie"`
		}{
			Message: "Movie deleted",
			Movie:   movie,
		})
	}
}

// updateMovieById function handles the PUT request to update a movie's details by ID
func updateMovieById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req) // Get the movie ID from URL
	var updatedMovie Movie
	_ = json.NewDecoder(req.Body).Decode(&updatedMovie) // Decode the updated movie from the request body

	found := false
	for idx, item := range movies {
		if item.ID == params["id"] {
			// Movie found, updating the details
			found = true
			item.Title = updatedMovie.Title
			item.Isbn = updatedMovie.Isbn
			item.Director = updatedMovie.Director
			movies[idx] = item // Save the updated movie back into the slice
			break
		}
	}

	// Return appropriate response
	w.Header().Set("Content-Type", "application/json")
	if found {
		// If movie is updated, return success message and updated movie details
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
			Movie   Movie  `json:"movie"`
		}{
			Message: "Movie updated",
			Movie:   updatedMovie,
		})
	} else {
		// If movie not found, return error message
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	}
}

// createMovie function handles the POST request to create a new movie
func createMovie(w http.ResponseWriter, req *http.Request) {
	var newMovie Movie
	_ = json.NewDecoder(req.Body).Decode(&newMovie) // Decode the new movie from the request body

	// Add the new movie to the movies slice
	movies = append(movies, newMovie)

	// Return response confirming the movie has been created
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Message: "Movie created",
		Movie:   newMovie,
	})
}

// getMovieById function handles the GET request to fetch a movie by ID
func getMovieById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req) // Get the movie ID from URL

	found := false
	var movie Movie
	for _, item := range movies {
		if item.ID == params["id"] {
			// Movie found, returning its details
			found = true
			movie = item
			break
		}
	}

	// Return appropriate response
	w.Header().Set("Content-Type", "application/json")
	if found {
		// If movie found, return movie details
		json.NewEncoder(w).Encode(movie)
	} else {
		// If movie not found, return error message
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Invalid Id",
		})
	}
}

// getMovies function handles the GET request to fetch all movies
func getMovies(w http.ResponseWriter, req *http.Request) {
	// Return the list of all movies
	fmt.Println("Fetching movies...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
