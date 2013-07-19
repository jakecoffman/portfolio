package handlers

import (
	"html/template"
	"net/http"
	"os"
	"portfolio/utils"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("base.html").Funcs(utils.FuncMap).ParseFiles(
		os.Getenv("WEBROOT")+"web/templates/base.html",
		os.Getenv("WEBROOT")+"web/templates/index.html",
	))

	session, _ := utils.Store.Get(r, "session")
	context := make(map[string]interface{})
	context["flashes"] = session.Flashes()
	session.Save(r, w)

	if err := t.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
