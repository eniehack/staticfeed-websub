package handler

import (
	"github.com/alexedwards/scs/v2"
	"github.com/eniehack/staticfeed-websub/internal/config"
	"github.com/eniehack/staticfeed-websub/internal/models"
	"golang.org/x/oauth2"
)

type Handler struct {
	DB           *models.Queries
	Config       *config.Config
	OAuth2Config *oauth2.Config
	Session      *scs.SessionManager
}
