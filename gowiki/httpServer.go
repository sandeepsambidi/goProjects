package main

import (
  "fmt"
  "log"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "url path received in req: %s", r.URL.Path[1:])
}

func main() {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(":8080",nil))
}
