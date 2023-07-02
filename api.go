package main

type ShortUrl struct {
	Status    string `json:"status"`
	Short_url string `json:"short_url"`
}

// func getShortUrls(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// }

func hello() {
	println("hello world")
}

func main() {

	hello()
	suii()
	// http.HandleFunc("/", getRoot)
	// fmt.Println("server starting at port 3000")
	// err := http.ListenAndServe(":3000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
