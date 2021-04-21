package main

import (
    "./CNAA/microservice"
    "log"
)

func main() {
    s := microservice.NewServer("", "8000")
    log.Fatal(s.ListenAndServe())
}
