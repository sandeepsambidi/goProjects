package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
  "log"
  "regexp"
  "errors"
  "path/filepath"
  "os"
)

type Page struct {
  Title string
  Body []byte
}

var (
  pwd, _ = os.Getwd()
  templatesPath = filepath.Join(pwd, "tmpl")
  dataPath = filepath.Join(pwd, "data")

  templates = template.Must(template.ParseFiles(filepath.Join(templatesPath,"edit.html"),
        filepath.Join(templatesPath,"view.html")))
  validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filepath.Join(dataPath,filename), p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, _ := ioutil.ReadFile(filepath.Join(dataPath,filename))
  return &Page{Title: title, Body: body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w,r)
      return
    }
    fn(w, r, m[2])
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html",p)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}

func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
