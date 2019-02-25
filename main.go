package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type point struct {
	x, y, z float64
}

// make a request to the nfw api and return the result
func netNFW(x float64, y float64, z float64) float64 {

	var nfwurl = os.Getenv("nfwurl")

	// build the request string
	var randMinRequestURL string = fmt.Sprintf("http://%s/NFW?x=%f&y=%f&z=%f", nfwurl, x, y, z)

	// make the request
	resp, err := http.Get(randMinRequestURL)
	if err != nil {
		panic(err)
	}

	// close the request body when finished
	defer resp.Body.Close()

	// read the body
	body, bodyerr := ioutil.ReadAll(resp.Body)
	if bodyerr != nil {
		panic(bodyerr)
	}

	var dat map[string]interface{}

	jsonerr := json.Unmarshal(body, &dat)
	if jsonerr != nil {
		panic(jsonerr)
	}

	result := dat["NFW"].(float64)

	return result
}

func gen(galaxyRange float64) point {

	// Define variables
	var length float64 = galaxyRange

	// define the range of the galaxy
	var rangeMin float64 = -length
	var rangeMax float64 = length

	var randMin float64 = netNFW(0, 0, 0)
	fmt.Printf("randmin: %30.20f", randMin)
	var randMax float64 = netNFW(length, length, length)
	fmt.Printf("randmax: %30.20f", randMax)

	var starFound bool = false

	for starFound == false {

		randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

		// generate random coordinates
		var x float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin
		var y float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin
		var z float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin

		var randomValue = randomSource.Float64()
		var randVal = ((randMax - randMin) * randomValue) + randMin

		if randVal < netNFW(x, y, z) {
			var newStar point = point{x, y, z}

			starFound = true
			return newStar
		}
	}
	return point{0, 0, 0}
}

func generate(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	amount, _ := strconv.ParseInt(params.Get("num"), 10, 64)
	galaxyRange, _ := strconv.ParseFloat(params.Get("range"), 64)

	// generate the given amount of stars
	for i := 0; i < int(amount); i++ {
		result := gen(galaxyRange)
		log.Printf("galaxy range: %f", galaxyRange)
		log.Printf("%v\n", result)
		_, err := fmt.Fprintf(w, "%f, %f, %f\n", result.x, result.y, result.z)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/gen", generate).Methods("GET")
	log.Fatal(http.ListenAndServe(":8123", router))
}
