package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

///update?item=socks&price=6 syntax
func main() {
	db := database{"shoes": 50, "socks": 5} //create an instance of type database
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

var mutex = &sync.RWMutex{}
type dollars float64 //declare type dollars
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) } //only keep 2 decimal places of dollars
type database map[string]dollars //databaseis a map of items and their dollar values


func (db database) list(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()                 //locks for reading
	fmt.Fprintf(w, "UPDATE HERE")
	for item, price := range db { //print items in the database
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	mutex.RUnlock()
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	mutex.RUnlock()
}

//Create Handler
func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	//Checking for error on price 
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Price: %v\n", err)
		return
	} else {
		mutex.Lock() 
		db[item] = dollars(p)
		mutex.Unlock()
		fmt.Fprint(w, "Item: ", item, ". Item price: ", p, "\n")
	}
}

//Update Handler
func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Price: %v\n", err)
		return
	}
	if _, found := db[item]; found {
		mutex.Lock()
		db[item] = dollars(p)
		mutex.Unlock()
		fmt.Fprint(w, "Updated item: ", item, ". Item price:  ", p, " \n") 
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, found := db[item]; found {
		mutex.Lock()
		delete(db, item)
		mutex.Unlock()
		fmt.Fprintf(w,"Deleted item %s\n", item)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}