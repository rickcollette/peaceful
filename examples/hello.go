package main

import (
    "fmt"
    "net/http"
    "github.com/rickcollette/peaceful/router"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(w, "Hello, world!")
}

func main() {
    router := restful.NewRouter()
    router.Handle("GET", "/hello", helloHandler)

    http.Handle("/", router)
    http.ListenAndServe(":8080", nil)
}
