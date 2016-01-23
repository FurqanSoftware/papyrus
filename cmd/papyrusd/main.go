package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gophergala2016/papyrus/ui"
)

func main() {
	http.Handle("/", ui.NewServer())

	log.Printf("Ligtening on %s", os.Getenv("ADDR"))
	err := http.ListenAndServe(os.Getenv("ADDR"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
