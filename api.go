package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Long_url_struct struct {
	Input_long_url string `json:"long_url"`
}

type Short_url_struct struct {
	Input_short_url string `json:"short_url"`
}

type response_obj_struct struct {
	Status    string `json:"status"`
	Short_url string `json:"short_url"`
	Long_url  string `json:"long_url"`
}

type db_document_struct struct {
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

	// check if wrong request method
	if r.Method != "POST" {
		http.Error(w, "only POST supported in encode method", http.StatusNotFound)
		return
	}

	// check if empty body, can be done via io.EOF as wel
	if r.Body == http.NoBody {
		http.Error(w, "empty body in POST in encode method", http.StatusNotAcceptable)
		return
	}
	// err := r.ParseForm() wont work as data not in www-x-formencoded but raw data

	decoder := json.NewDecoder(r.Body) // it buffers the entire json value in memory before unmarshal
	var long_url_obj Long_url_struct
	err := decoder.Decode(&long_url_obj) // unmarshal/decoding occurs here
	if err != nil {
		panic(err)
	}

	// mongo filter to find from given long url
	filter := bson.D{{"long_url", long_url_obj.Input_long_url}}

	doc_count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var res_obj response_obj_struct

	// check if already exists
	if doc_count > 0 {
		// it already exists
		var db_obj db_document_struct
		fmt.Println("document of long url already exists")
		err = collection.FindOne(context.TODO(), filter).Decode(&db_obj)

		if err != nil {
			panic(err)
		}

		res_obj = response_obj_struct{
			Status:    "existed",
			Short_url: db_obj.Short_url,
			Long_url:  db_obj.Long_url,
		}

	} else {
		new_short_url := getTinyUrl()

		// insert document in DB for same
		insert_doc := db_document_struct{Short_url: new_short_url, Long_url: long_url_obj.Input_long_url}
		insert_res, err := collection.InsertOne(context.TODO(), insert_doc)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("document inserted of id: ", insert_res.InsertedID)

		res_obj = response_obj_struct{
			Status:    "created",
			Short_url: insert_doc.Short_url,
			Long_url:  insert_doc.Long_url,
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

// decode the long_url from short url
func decode(w http.ResponseWriter, r *http.Request) {

	// check if wrong request method
	if r.Method != "POST" {
		http.Error(w, "only POST method accepted in decode method", http.StatusNotFound)
	}

	// check if empty body, can be done via io.EOF as well
	if r.Body == http.NoBody {
		http.Error(w, "empty body in POST in decode method", http.StatusNotAcceptable)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var short_url_obj Short_url_struct
	err := decoder.Decode(&short_url_obj) // unmarshal/decode by reference
	if err != nil {
		panic(err)
	}

	// mongodb filter to check if short url in DB
	filter := bson.D{{"short_url", short_url_obj.Input_short_url}}

	doc_count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var res_obj response_obj_struct

	if doc_count > 0 {
		// it exists in DB
		var db_obj db_document_struct
		fmt.Println("document of short url already exists")
		collection.FindOne(context.TODO(), filter).Decode(&db_obj) // decode in object of given struct

		res_obj = response_obj_struct{
			Status:    "existed",
			Short_url: db_obj.Short_url,
			Long_url:  db_obj.Long_url,
		}

	} else {
		println("long url does not exist")
		// it just does not exist
		res_obj = response_obj_struct{
			Status:    "not_found",
			Short_url: short_url_obj.Input_short_url, // dont use db_obj url here as it is empty
			Long_url:  "",
		}
	}

	// encode the res_obj in json
	jsonData, err := json.Marshal(res_obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func getRoot(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "only GET method supported at /", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "hello suii")

}

func apiMethod() {
	router := mux.NewRouter()
	router.HandleFunc("/", getRoot)
	router.HandleFunc("/encode", encode)
	router.HandleFunc("/decode", decode)
	http.Handle("/", router) // important as here we define the router

	fmt.Printf("server started at port 9000")
	err := http.ListenAndServe(":9000", nil)
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
