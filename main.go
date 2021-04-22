package main

import (
    "src/microservice"
    "log"
)

func main() {
    s := microservice.NewServer("", "8000")
    log.Fatal(s.ListenAndServe())
}
