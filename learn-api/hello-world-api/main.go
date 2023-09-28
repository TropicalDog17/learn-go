package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// marshall, unmarshall
// marshall: go obj -> json
// unmarshall: json -> go
// {"name": "tuan", "age": 12, "description": "Very handsome"}
// == Person{Name: "tuan", Age: 12, Desc: "Very handsome"}
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Desc string `json:"description"`
}

// Global variable
var people []Person

// entrypoint of API
func home(w http.ResponseWriter, r *http.Request) {
	// Response
	fmt.Fprintf(w, "Welcome to the homepage")
	fmt.Println("Endpoint hit: Homepage")
}
func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllPeople")
	// Encode json to write to Response
	json.NewEncoder(w).Encode(people)
}

// Url: localhost:3000/person?id=1
func getByIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: getByIndex")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	peopleCount := len(people)
	if peopleCount == 0 {
		fmt.Println("Empty people list")
		json.NewEncoder(w).Encode([]Person{})
		return
	}
	if int(id) >= peopleCount || int(id) < 0 {
		fmt.Println("Index out of range")
		json.NewEncoder(w).Encode([]Person{})
		return
	}

	person := &people[int(id)]
	json.NewEncoder(w).Encode(person)
}
func addPerson(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newPerson Person
	// TODO: handle if reqBody JSON doesn't match the Person schema.
	json.Unmarshal(reqBody, &newPerson)

	people = append(people, newPerson)
	// Return the success response
	resp := &Response{}
	resp.Message = "Add person successfully!"
	resp.Success = true
	resp.Data = newPerson
	json.NewEncoder(w).Encode(resp)
}
func updatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	resp := &Response{}
	if !isValidIndex(id) {
		resp.Success = false
		resp.Message = "Invalid index :("
		resp.Data = nil
		json.NewEncoder(w).Encode(resp)
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updatedPerson Person
	json.Unmarshal(reqBody, &updatedPerson)

	people[id] = updatedPerson
	resp.Success = true
	resp.Message = fmt.Sprintf("Successfully update person index %v", id)
	resp.Data = updatedPerson
	json.NewEncoder(w).Encode(resp)
}
func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
func deletePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint DeletePerson")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	resp := &Response{}
	if !isValidIndex(id) {
		resp.Success = false
		resp.Message = "Invalid index :("
		resp.Data = nil
		json.NewEncoder(w).Encode(resp)

		return
	}
	deletedPerson := people[id]
	resp.Data = deletedPerson
	resp.Message = fmt.Sprintf("Successfully removed the person at id %v. Keep in mind that all other person id is changed :)", id)
	resp.Success = true
	people = RemoveIndex[Person](people, id)
	json.NewEncoder(w).Encode(resp)
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/all", getAll)
	myRouter.HandleFunc("/person/{id}", getByIndex).Methods("GET")
	myRouter.HandleFunc("/person/{id}", updatePerson).Methods("PUT")
	myRouter.HandleFunc("/person/{id}", deletePerson).Methods("DELETE")
	myRouter.HandleFunc("/person", addPerson).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}
func isValidIndex(id int) bool {
	if id < 0 {
		return false
	}
	if id >= len(people) {
		return false
	}
	return true
}
func main() {
	people = []Person{
		{Name: "tuan", Age: 12, Desc: "A very handsome gigachad"},
		{Name: "wojak", Age: 25, Desc: "A loser in life"},
	}
	handleRequests()
}
