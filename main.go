package main

import (
	"fmt"
	"log"
	"net/http"
	"restore-photos/router"
)

func main() {
	fmt.Println("Backend for RestorePhotos")

	r := router.Router()

	err := http.ListenAndServe(":4068", r)
	if err != nil {
		log.Fatal(err)
	}
}
