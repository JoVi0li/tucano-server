package main

import (
	"log"
	"net/http"

	"github.com/JoVi0li/tucano-server/internal"
)

func main() {
	clientList := internal.NewClientList()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleSignalling(w, r, clientList)
	})

	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatal(err)
	}

}
