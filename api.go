package main

type ShortUrl struct {
	Status    string `json:"status"`
	Short_url string `json:"short_url"`
}

// func getShortUrls(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// }

func main() {

	apiMethod()
}
