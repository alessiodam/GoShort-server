package router

import (
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, nil)
}
