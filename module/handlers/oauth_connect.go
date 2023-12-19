package handlers

import (
	"context"
	"fmt"
	"github.com/jaevor/go-nanoid"
	"golang.org/x/oauth2"
	"net/url"
)

type OauthConnectRepo interface {
	GetOauthConfig() (*oauth2.Config, error)
}

type oauthConnectHdl struct {
	repo OauthConnectRepo
}

func NewOauthConnectHdl(repo OauthConnectRepo) *oauthConnectHdl {
	return &oauthConnectHdl{repo: repo}
}

func (h *oauthConnectHdl) Response(ctx context.Context, connectorId string) (string, error) {
	canonicID, _ := nanoid.Standard(21)
	state := canonicID()

	oauthConfig, err := h.repo.GetOauthConfig()
	if err != nil {
		return "", err
	}
	authUrl := oauthConfig.AuthCodeURL(state)

	u, err := url.Parse(authUrl)
	if err != nil {
		return "", err
	}

	if connectorId != "" {
		u.Path += fmt.Sprintf("/%s", connectorId)
	}

	return u.String(), nil
}
