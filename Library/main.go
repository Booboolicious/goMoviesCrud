package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	. "github.com/Booboolicious/my"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Movie struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Isbn      string    `json:"isbn"`
	Title     string    `json:"title"`
	Director  Director  `json:"director" gorm:"embedded;embeddedPrefix:director_"`
	Year      string    `json:"year"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("movies.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Movie{})

	// Seed data if empty
	var count int64
	db.Model(&Movie{}).Count(&count)
	if count == 0 {
		db.Create(&Movie{
			ID: "1",
			Isbn: "123456789",
			Title: "The Matrix",
			Director: Director{Firstname: "John", Lastname: "Doe"},
			Year: "1999",
		})
		db.Create(&Movie{
			ID: "2",
			Isbn: "124486789",
			Title: "Now You See Me",
			Director: Director{Firstname: "Louis", Lastname: "Leterrier"},
			Year: "2013",
		})
	}

	r := mux.NewRouter()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Serve Static Files for Frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

	Log("Server started at port 8000")
	http.ListenAndServe(":8000", corsMiddleware(r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movies []Movie
	db.Find(&movies)
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	result := db.First(&movie, "id = ?", params["id"])
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	db.Create(&movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	result := db.First(&movie, "id = ?", params["id"])
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var updatedMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)
	updatedMovie.ID = params["id"]
	db.Save(&updatedMovie)
	json.NewEncoder(w).Encode(updatedMovie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	db.Delete(&Movie{}, "id = ?", params["id"])
	
	// Return updated list to match old behavior
	var movies []Movie
	db.Find(&movies)
	json.NewEncoder(w).Encode(movies)
}