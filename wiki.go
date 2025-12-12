package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string

	Body []byte
}

func (p *Page) save() error {

	filename := "./pages/" + p.Title + ".txt"

	return os.WriteFile(filename, p.Body, 0600)

}

func loadPage(title string) (*Page, error) {

	filename := "./pages/" + title + ".txt"

	body, err := os.ReadFile(filename)

	if err != nil {

		return nil, err

	}

	return &Page{Title: title, Body: body}, nil

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".html")

	if err != nil {
		http.Error(w, "Failed to parse html", http.StatusBadGateway)
		return
	}

	t.Execute(w, p)
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./pages")

	if err != nil {
		http.Error(w, "Failed to read dir", http.StatusBadGateway)
		return
	}

	var pages []string
	for _, file := range files {
		// Remove .txt extension
		name := strings.TrimSuffix(file.Name(), ".txt")
		pages = append(pages, name)
	}

	t, err := template.ParseFiles("./templates/landing.html")
	if err != nil {
		http.Error(w, "Failed to parse html", http.StatusBadGateway)
		return
	}

	t.Execute(w, pages)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
