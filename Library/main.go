package main

import (
	."github.com/Booboolicious/my"
	"log"
	"encoding/json"
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
) 

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
	Year string `json:"year"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	r.
	Log("Hello World")
}