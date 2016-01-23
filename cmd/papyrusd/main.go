package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gophergala2016/papyrus/ui"
)

func main() {
	err := data.OpenDBSession(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = data.MakeIndexes()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", ui.NewServer())

	log.Printf("Ligtening on %s", os.Getenv("ADDR"))
	err = http.ListenAndServe(os.Getenv("ADDR"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
