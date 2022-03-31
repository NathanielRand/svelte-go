package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/NathanielRand/go-svelte/boilerplate/views"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const port = ":8080"

// Global variables related to templates.
// ATTN: Relocate for production. Global variables are bad practice.
var (
	indexView       *views.View
	homeView        *views.View
	notFound404View *views.View
)

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(0, 0, 1).Format(http.TimeFormat)
)

// Response type hold data from the api response
type Response struct {
	Data DataType `json:"data"`
}

type DataType struct {
	// STRUCTURE API DATA HERE
}

func makeAPIRequest(url string) []byte {
	apiURL := "https://exampleapiendpoint.com/value=" + url

	req, _ := http.NewRequest("GET", apiURL, nil)

	req.Header.Add("x-rapidapi-key", "")
	req.Header.Add("x-rapidapi-host", "")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "max-age:15780000, public")
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)

	switch r.Method {
	case "GET":
		must(homeView.Render(w, nil))
	case "POST":
		// Form data
		r.ParseForm()

		// Get field form input value
		url := r.FormValue("formValue")

		// Make the api request using the provided URL.
		APIResponseBody := makeAPIRequest(url)

		var responseObject Response
		json.Unmarshal(APIResponseBody, &responseObject)

		// Store gathered data in a object and assign to var
		results := Response{
			Data: responseObject.Data,
		}

		// Render home page with passed in data
		must(homeView.Render(w, results))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "max-age:15780000, public")
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)

	must(indexView.Render(w, nil))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

var notfound http.Handler = http.HandlerFunc(notFound404)

// notFound404 prints message to screen. TODO: replace with custom 404 page.
func notFound404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "max-age:15780000, public")
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)
	must(notFound404View.Render(w, nil))
}

// must is a helper for errors
func must(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {
	// Load .env file from given path
	// we keep it empty it will load .env from current directory
	// ATTN: Consdier moving this load to an init() func to be run before main() func.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	indexView = views.NewView("bulma", "views/pages/index.html")
	homeView = views.NewView("bulma", "views/pages/home.html")
	notFound404View = views.NewView("bulma", "views/pages/notFound404.html")

	// Gorilla Mux router
	r := mux.NewRouter()

	// Handle 404s
	r.NotFoundHandler = notfound

	// Assest Routes
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Public Routes
	publicHandler := http.FileServer(http.Dir("./static/"))
	publicHandler = http.StripPrefix("/static/", publicHandler)
	r.PathPrefix("/public/").Handler(publicHandler)

	// Favicon
	r.HandleFunc("/favicon.ico", faviconHandler)

	// Routes
	r.HandleFunc("/", index)
	r.HandleFunc("/home", home)

	// Start web server.
	fmt.Println("Listening on port", port)
	http.ListenAndServe(port, r)
}
