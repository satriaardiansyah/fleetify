package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server Fleetify backend berjalan di port 8080...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Halo! Ini adalah backend Fleetify.")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
