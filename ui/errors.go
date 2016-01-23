package ui

import "net/http"

func ServeNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func ServeUnauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func ServeBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func ServeInternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func catch(r *http.Request, err error) {
	if err != nil {
		panic(err)
	}
}
