package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func main() {
	server("localhost:5005")
}

func server(listen string) {
	m := pat.New()
	m.Post("/locks", http.HandlerFunc(create))
	m.Get("/locks/:id", http.HandlerFunc(show))
	m.Get("/locks", http.HandlerFunc(index))
	m.Del("/locks/:id", http.HandlerFunc(del))

	http.Handle("/", m)

	log.Printf("Starting server on %s\n", listen)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func create(w http.ResponseWriter, req *http.Request) {
	params := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	lock, acquired := reg.set(params["id"])
	data, err := json.Marshal(lock)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if acquired {
		log.Println("create", http.StatusCreated)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Println("create", http.StatusLocked)
		w.WriteHeader(http.StatusLocked)
	}

	fmt.Fprintf(w, string(data))
}

func show(w http.ResponseWriter, req *http.Request) {
	log.Println("show")

	id := req.URL.Query().Get(":id")

	lock, exist := reg.get(id)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(lock)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))
}

func index(w http.ResponseWriter, req *http.Request) {
	log.Println("index")

	locks := reg.all()

	data, err := json.Marshal(locks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))
}

func del(w http.ResponseWriter, req *http.Request) {
	log.Println("delete")

	id := req.URL.Query().Get(":id")

	if !reg.unset(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
