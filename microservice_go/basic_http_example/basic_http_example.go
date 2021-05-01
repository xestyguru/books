package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
type validationContextKey string

type helloWorldResponse struct {
	Message string `json:"message"`
	Author  string `json:"-"`
	Date    string `json:", omitempty"`
	Id      int    `json:"id, string"`
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

func main() {
	port := 8080

	cathandler := http.FileServer(http.Dir("./images"))
	handler := newValidationHandler(newHelloWorldHandler())

	http.Handle("/helloworld", handler)
	http.Handle("/cat/", http.StripPrefix("/cat/",cathandler))
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

type validationsHandler struct {
	next http.Handler
}

func newValidationHandler(next http.Handler) http.Handler {
	return validationsHandler{next: next}
}

func (h validationsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	c := context.WithValue(r.Context(), validationContextKey("name"), request.Name)
	r = r.WithContext(c)

	h.next.ServeHTTP(rw, r)
}

type helloWorldHandler struct {}

func newHelloWorldHandler() http.Handler{
	return helloWorldHandler{}
}

func (h helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(validationContextKey("name")).(string)
	response := helloWorldResponse{Message: "Hello " + name}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

// func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	// data, err := json.Marshal(response)
// 	// if err != nil {
// 	// 	panic("Ooops")
// 	// }

// 	// fmt.Fprint(w, string(data))

// 	// body, err := ioutil.ReadAll(r.Body)
// 	// if err != nil {
// 	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	var request helloWorldRequest
// 	// err = json.Unmarshal(body, &request)
// 	// if err != nil {
// 	// 	http.Error(w, "Bad request2", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	decoder := json.NewDecoder(r.Body)

// 	err := decoder.Decode(&request)
// 	if err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	response := helloWorldResponse{Message: "Hello " + request.Name}

// 	encoder := json.NewEncoder(w)
// 	encoder.Encode(&response)
// }
