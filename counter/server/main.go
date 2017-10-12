package main

import (
    "fmt"
    "net/http"
	"simplesurance-group.de/counter/storage"
)


const (
    countersFilePath = "counters.log"
)

func indexHandler() func(http.ResponseWriter, *http.Request) {
	store := storage.NewTimestampStorage(countersFilePath)
    return func(w http.ResponseWriter, r *http.Request) {
        c_out := store.CounterAddTimestampNow()
        counter := <- c_out
        fmt.Fprintf(w, "%d", counter)
    }
}

func main() {
    http.HandleFunc("/", http.HandlerFunc(indexHandler()))
    http.ListenAndServe(":8080", nil)
}
