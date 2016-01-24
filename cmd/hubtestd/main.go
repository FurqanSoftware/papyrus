package main

import (
	"net/http"

	"github.com/desertbit/glue"
	"github.com/gophergala2016/papyrus/hub"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/hubtestd/index.html")
}

func main() {
	http.HandleFunc("/", ServeIndex)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("cmd/hubtestd/assets"))))

	hub.DefaultRepository.GetBlob = func(id string) ([]byte, error) {
		if id != "1" {
			return nil, nil
		}
		return []byte{}, nil
	}

	gs := glue.NewServer()
	gs.OnNewSocket(hub.HandleSocket)
	http.Handle("/glue/", gs)

	err := http.ListenAndServe(":8080", nil)
	catch(err)
}
