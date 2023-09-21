package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/eniehack/staticfeed-websub/internal/config"
	"github.com/eniehack/staticfeed-websub/internal/handler"
	"github.com/eniehack/staticfeed-websub/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	_ "modernc.org/sqlite"
)

func main() {
	h := new(handler.Handler)
	config, err := config.LoadFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("config.LoadFromFile error: %v", err)
		os.Exit(1)
	}
	h.Config = config
	db, err := sql.Open("sqlite", config.DataBaseConfig.Url)
	if err != nil {
		log.Fatalf("db: %v", err)
		os.Exit(1)
	}
	h.DB = models.New(db)
	if h.DB == nil {
		log.Fatalf("h.DB is nil")
	}

	h.OAuth2Config = &oauth2.Config{
		ClientID:     config.GitHubConfig.ClientID,
		ClientSecret: config.GitHubConfig.ClientSecret,
		Scopes:       []string{"user:read"},
		Endpoint:     github.Endpoint,
	}
	h.Session = scs.New()
	h.Session.Lifetime = 24 * time.Hour

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/", h.GetTopPage)
	r.Get("/signin/github", h.GitHubLogin)
	r.Get("/signin/github/callback", h.CallbackFromGitHub)

	log.Println("Listening on port:3000")

	log.Fatal(
		http.ListenAndServe(":3000", h.Session.LoadAndSave(r)),
	)
}
