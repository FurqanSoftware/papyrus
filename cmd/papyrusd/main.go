package main

import (
	"log"
	"net/http"
	"os"

	"github.com/desertbit/glue"
	"github.com/gophergala2016/papyrus/data"
	"github.com/gophergala2016/papyrus/hub"
	"github.com/gophergala2016/papyrus/repo"
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

	hub := hub.New(repo.New())

	gs := glue.NewServer()
	gs.OnNewSocket(hub.HandleSocket)
	http.Handle("/glue/", gs)

	http.Handle("/", ui.NewServer())

	log.Printf("Listening on %s", os.Getenv("ADDR"))
	err = http.ListenAndServe(os.Getenv("ADDR"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	loadHerokuEnv()
}
