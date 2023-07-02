package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// init method of golang called only once at start
// func init() {
// 	seed := time.Now().UnixNano()
// 	rand.New(rand.NewSource(seed))
// }

// here we are defining and not declaring
var chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var short_to_long_mapper map[string]string // key is string and value is int, make is to make empty map
var long_to_short_mapper map[string]string

type Long_url_struct struct {
	Long_url string `json:"long_url"`
}

type Short_url_struct struct {
	Short_url string `json:"short_url"`
}

type response_struct struct {
	Status    string `json:"status"`
	Short_url string `json:"short_url"`
	Long_url  string `json:"long_url"`
}

// generate random tiny-url for length 7
func getTinyUrl() string {
	short_url := ""
	for i := 0; i < 7; i++ {
		r := rand.Intn(len(chars) - 1)
		short_url += string(chars[r])
	}
	short_url = "http://tiny-url.com/" + short_url
	fmt.Println(short_url)
	return short_url
}

func encode(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm() wont work as data not in www-x-formencoded but raw data

	decoder := json.NewDecoder(r.Body) // it buffers the entire json value in memory before unmarshal
	var long_url_obj Long_url_struct
	err := decoder.Decode(&long_url_obj) // unmarshal/decoding occurs here
	if err != nil {
		panic(err)
	}
	// log.Println(long_url_obj.Long_url)

	// w.WriteHeader(http.StatusOK)
	short_url, ok := long_to_short_mapper[long_url_obj.Long_url] // ok is bool, we can instead use if mapper[url]{} as well instead

	// empty struct with no value or data assigned rn
	var res_obj response_struct
	if ok {
		res_obj = response_struct{
			Status:    "existed",
			Short_url: short_url,
			Long_url:  long_url_obj.Long_url,
		}
	} else {
		new_short_url := getTinyUrl()
		long_to_short_mapper[long_url_obj.Long_url] = short_url // store the value
		res_obj = response_struct{
			Status:    "created",
			Short_url: new_short_url,
			Long_url:  long_url_obj.Long_url,
		}
	}

	// encode it to json before sending
	jsonData, err := json.Marshal(res_obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// todo: WIP
func decode(w http.ResponseWriter, r *http.Request) {

}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello suii")

}

func suii() {
	router := mux.NewRouter()
	short_to_long_mapper = make(map[string]string)
	long_to_short_mapper = make(map[string]string)
	// new_url := "https://github.com"
	// http.Handle("/", http.RedirectHandler(new_url, http.StatusMovedPermanently))

	router.HandleFunc("/", getRoot).Methods("GET")
	router.HandleFunc("/encode", encode).Methods("POST")
	http.Handle("/", router) // important as here we define the router

	fmt.Printf("server started at port 3000")
	// log.Fatal(http.ListenAndServe(":3000", nil))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("404 error bruh")
		log.Fatal(err)
	}
}

// first we try without using mongoDB
// use request.ParseForm() only when header of type `application/x-www-form-urlencoded/`
// otherwise use ioutil.ReadAll()

// we can either read the body using ioutil.ReadAll() and then json.UnMarshal
// bodyBuffer, err := ioutil.ReadAll(r.Body) // gives raw data in bytes and error if any
// if err != nil {
// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 	return // no further addition
// }
// and then map to struct long_url_struct using unMarshal()

// NOTE: json.NewDecoder() is better as used below and in code

// or use decoder: = json.NewDecoder(request.body)
// and then store the map the json to a struct. Read more here
// https://articles.wesionary.team/difference-of-json-encoding-vs-marshaling-and-json-decoding-vs-unmarshaling-1a6baf6a7f5c

// Sprint is another print method useful
// res_str = fmt.Sprint("short url ", short_url, "already exists")
