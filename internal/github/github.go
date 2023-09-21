package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	UserData struct {
		Login string `json:"login"`
		Name  string `json:"name"`
	}
)

var ErrForbiddenStatus error = errors.New("github API returns Forbidden(403).")

func GetUserData(ctx context.Context, token string) (*UserData, error) {
	client := new(http.Client)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrForbiddenStatus
	}

	user := new(UserData)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}
