package main

import (
  "sync"
  "math/rand"
  "time"
  "fmt"
  "net/http"
  "html/template"
  )


const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`

var src = rand.NewSource(time.Now().UnixNano())
var store = NewURLStore()
//var templates = template.Must(template.ParseFiles("edit.html"))

type URLStore struct {
  urls map[string]string
  mu sync.RWMutex
}

func (urlStore *URLStore) Get(key string) string {
  urlStore.mu.RLock()
  defer urlStore.mu.RUnlock()
  return urlStore.urls[key]
}

func (urlStore *URLStore) Set(key, url string) bool {
  urlStore.mu.Lock()
  defer urlStore.mu.Unlock()
  _, exists := urlStore.urls[key]

  if exists {
    return false
  }
  urlStore.urls[key] = url
  return true
}

func (urlStore *URLStore) Put(url string) string {
  for {
    key := genKey(5)
    if urlStore.Set(key, url) {
      return key
    }
  }
  panic("could not set key for the url ")
}

func genKey(n int) string {
  b := make([]byte, n)
   // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
   for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
       if remain == 0 {
           cache, remain = src.Int63(), letterIdxMax
       }
       if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
           b[i] = letterBytes[idx]
           i--
       }
       cache >>= letterIdxBits
       remain--
   }

   return string(b)
}

func NewURLStore() *URLStore {
  u := &URLStore{urls:make(map[string]string)}
  return u;
}

func addHandler(w http.ResponseWriter, r *http.Request) {
  url := r.FormValue("url")
  if url == "" {
    //fmt.Fprintf(w, AddForm)
    // fmt.Fprintf(w, "<h1>Add Url</h1>"+
    //     "<form action=\"/add\" method=\"POST\">"+
    //     "<input type=\"text\" name=\"url\">"+
    //     "<input type=\"submit\" value=\"Add\">"+
    //     "</form>")
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, nil)
    return
  }
  key := store.Put(url)
  fmt.Fprintf(w, "http://localhost:8080/%s", key)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
  key:= r.URL.Path[1:]
  url := store.Get(key)

  if url == "" {
    http.NotFound(w,r)
    return
  }
  http.Redirect(w, r, url, http.StatusFound)
}

func main() {
  http.HandleFunc("/add",addHandler)
  http.HandleFunc("/", redirectHandler)
  http.ListenAndServe(":8080", nil)
}
