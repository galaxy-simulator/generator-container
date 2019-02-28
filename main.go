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
func netNFW(x float64, y float64, z float64) (float64, error) {

	// get the url of the nfw container or a reverse proxy handling the request to the containers
	var nfwurl = os.Getenv("nfwurl")

	// if no url is given, ask
	if nfwurl == "" {
		return 0, fmt.Errorf("No nfwurl given!")
	}

	// build the request string
	var randMinRequestURL string = fmt.Sprintf("http://%s/NFW?x=%f&y=%f&z=%f", nfwurl, x, y, z)

	// make the request
	resp, err := http.Get(randMinRequestURL)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	// read the body
	body, bodyerr := ioutil.ReadAll(resp.Body)
	if bodyerr != nil {
		panic(bodyerr)
	}

	// define an interface into which the nfw gets unpacked
	var dat map[string]interface{}
	jsonerr := json.Unmarshal(body, &dat)
	if jsonerr != nil {
		panic(jsonerr)
	}

	result := dat["NFW"].(float64)

	return result, nil
}

func gen(galaxyRange float64) point {

	// Define variables
	var length float64 = galaxyRange

	// define the range of the galaxy
	var rangeMin float64 = -length
	var rangeMax float64 = length

	// get the minimal NFW value
	var randMin, errGetMinValue = netNFW(0, 0, 0)
	if errGetMinValue != nil {
		panic(errGetMinValue)
	}

	// get the maximal NFW value
	var randMax, errGetMaxValue = netNFW(length, length, length)
	if errGetMaxValue != nil {
		panic(errGetMaxValue)
	}

	var starFound bool = false

	for starFound == false {

		// define a new random source (without this, the numbers would not be random!)
		randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

		// generate random coordinates
		var x float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin
		var y float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin
		var z float64 = ((rangeMax - rangeMin) * randomSource.Float64()) + rangeMin

		// generate a random value in the (randmin, randmax) range
		var randomValue = randomSource.Float64()
		var randVal = ((randMax - randMin) * randomValue) + randMin

		// calculate the nfw-value of the previously generated star
		var nfwVal, err = netNFW(x, y, z)
		if err != nil {
			panic(err)
		}

		// check if th star should be kept or not
		if randVal < nfwVal {
			var newStar point = point{x, y, z}
			starFound = true
			return newStar
		}
	}

	// if no star is found at all, return (0, 0, 0)
	// this code should actually never be reached
	return point{0, 0, 0}
}

// the generate handler gets a number of stars and a range in which the stars should be generated
// and generated them
func generateHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// parse the url arguments
	amount, _ := strconv.ParseInt(params.Get("num"), 10, 64)
	galaxyRange, _ := strconv.ParseFloat(params.Get("range"), 64)

	// generate the given amount of stars
	for i := 0; i < int(amount); i++ {

		// generate the star
		result := gen(galaxyRange)

		log.Printf("galaxy range: %f", galaxyRange)
		log.Printf("%v\n", result)

		log.SetOutput(w)
		log.Printf("%f, %f, %f\n", result.x, result.y, result.z)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Generator container here!\nUse the /gen?num=<num>&range=<range> endpoint to generator a star!")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/gen", generateHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8123", router))
}
