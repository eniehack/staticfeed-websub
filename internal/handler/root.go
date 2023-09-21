package handler

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/oauth2/github"
)

func (h *Handler) GetTopPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./template/index.tmpl.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, map[string]string{"github_client_id": h.Config.GitHubConfig.ClientID, "github_oauth_url": github.Endpoint.TokenURL}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
