package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var StatusMap = make(map[string]string)

func main() {
	fmt.Println("Server is running")
	http.HandleFunc("/websites", handler)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Printf("Server error %v", err)
		return
	}
}

func checker(website string, task chan string) {
	if _, err := http.Get(website); err != nil {
		StatusMap[website] = "DOWN"
	} else {
		StatusMap[website] = "UP"
	}
	task <- website
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		param := r.URL.Query().Get("name")
		if param == "" {
			for key, val := range StatusMap {
				fmt.Fprintln(w, key, " ", val)
			}
		} else {
			status := StatusMap[param]
			fmt.Fprintln(w, fmt.Sprintf("Status of %s: %s", param, status))
		}
	} else if r.Method == http.MethodPost {
		var websites []string
		err := json.NewDecoder(r.Body).Decode(&websites)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Error in post method: %+v", err))
		}
		task := make(chan string)
		for _, url := range websites {
			go checker(url, task)
		}
		for website := range task {
			go func(web string) {
				time.Sleep(time.Minute)
				checker(web, task)
			}(website)
		}
	} else {
		fmt.Fprint(w, "Invalid request")
	}
}
