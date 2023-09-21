package handler

import (
	"html/template"
	"net/http"

	"github.com/eniehack/staticfeed-websub/internal/github"
)

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	id := h.Session.GetString(r.Context(), "user_id")
	user, err := h.DB.GetUserFromID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userGitHubData, err := github.GetUserData(r.Context(), user.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template.New("./template/dashboard/index.tmpl.html")
	if err := tmpl.Execute(w, map[string]string{"username": userGitHubData.Name}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
