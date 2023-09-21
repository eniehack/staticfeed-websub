package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/eniehack/staticfeed-websub/internal/models"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
)

func (h *Handler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.OAuth2Config.AuthCodeURL(generateStateOAuthCookie(w)), http.StatusTemporaryRedirect)
}

func (h *Handler) CallbackFromGitHub(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	url_state := r.FormValue("state")
	cookie_state, err := r.Cookie("oauth_state")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(code) <= 0 {
		log.Fatalf("callbackfromgithub: 400 (code is empty)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url_state != cookie_state.Value {
		log.Printf("OAuth state unmatched")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.OAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		log.Fatalf("callbackfromgithub: 400 (token is empty, %v)", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := h.InsertNewToken(r.Context(), token)
	if err != nil {
		log.Fatalf("InsertNewToken: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Session.Put(r.Context(), "id", id)
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func (h *Handler) InsertNewToken(ctx context.Context, token *oauth2.Token) (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), ulid.DefaultEntropy())
	if err != nil {
		return "", err
	}

	params := models.InsertNewTokenParams{
		AccessToken:  token.AccessToken,
		ID:           id.String(),
		RefreshToken: token.RefreshToken,
		ExpiredAt:    token.Expiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.DB.InsertNewToken(ctx, params); err != nil {
		return "", nil
	}
	return id.String(), nil
}

type MyTokenSource struct {
	Handler *Handler
	Src     oauth2.TokenSource
	DBID    string
}

func (h *Handler) NewClient(ctx context.Context, token *oauth2.Token) *http.Client {
	toksrc := h.OAuth2Config.TokenSource(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	src := &MyTokenSource{
		Handler: h,
		Src:     toksrc,
		DBID:    "",
	}
	reusesrc := oauth2.ReuseTokenSourceWithExpiry(token, src, time.Hour*8)
	return oauth2.NewClient(ctx, reusesrc)
}

func (h *Handler) UpdateToken(token *oauth2.Token, DBID string) error {
	ctx := context.Background()
	if err := h.DB.UpdateToken(ctx, models.UpdateTokenParams{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiredAt:    token.Expiry,
		UpdatedAt:    time.Now(),
		ID:           DBID,
	}); err != nil {
		return err
	}
	return nil
}

func (s *MyTokenSource) Token() (*oauth2.Token, error) {
	t, err := s.Src.Token()
	if err != nil {
		return nil, err
	}
	if err := s.Handler.UpdateToken(t, s.DBID); err != nil {
		return t, err
	}

	return t, nil
}

func generateStateOAuthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(20 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return state
}
