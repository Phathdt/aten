package handlers

import (
	"aten/module/models"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"aten/shared/errorx"
	"context"
	"errors"
	"github.com/coreos/go-oidc"
	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/oauth2"
)

type OauthCallbackRepo interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
	CreateUser(ctx context.Context, data *models.UserCreate) error
}

type OauthCallbackSessionRepo interface {
	SetUserToken(ctx context.Context, userId int, token, subToken string, expiredTime int) error
}

type OauthCallbackDexRepo interface {
	GetOauthConfig() (*oauth2.Config, error)
	GetIdTokenVerifier() (*oidc.IDTokenVerifier, error)
}

type oauthCallbackHdl struct {
	repo          OauthCallbackRepo
	sRepo         OauthCallbackSessionRepo
	dexRepo       OauthCallbackDexRepo
	tokenProvider tokenprovider.Provider
}

func NewOauthCallbackHdl(repo OauthCallbackRepo, sRepo OauthCallbackSessionRepo, dexRepo OauthCallbackDexRepo, tokenProvider tokenprovider.Provider) *oauthCallbackHdl {
	return &oauthCallbackHdl{repo: repo, sRepo: sRepo, dexRepo: dexRepo, tokenProvider: tokenProvider}
}

func (h *oauthCallbackHdl) Response(ctx context.Context, code string) (tokenprovider.Token, error) {
	oauthConfig, err := h.dexRepo.GetOauthConfig()
	if err != nil {
		return nil, err
	}

	oauth2Token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing raw id token")
	}

	idTokenVerifier, err := h.dexRepo.GetIdTokenVerifier()
	if err != nil {
		return nil, err
	}
	// Parse and verify ID Token payload.
	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	// Extract custom claims.
	var claims struct {
		Email           string `json:"email"`
		Name            string `json:"name"`
		FederatedClaims struct {
			ConnectorID string `json:"connector_id"`
			UserID      string `json:"user_id"`
		} `json:"federated_claims"`
	}
	if err = idToken.Claims(&claims); err != nil {
		return nil, err
	}

	user, err := h.repo.GetUserByCondition(ctx, map[string]interface{}{"email": claims.Email})
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		canonicID, _ := nanoid.Standard(21)
		data := models.UserCreate{
			SQLModel: core.NewSQLModel(),
			Email:    claims.Email,
			Password: canonicID(),
		}

		if err = h.repo.CreateUser(ctx, &data); err != nil {
			return nil, core.ErrBadRequest.
				WithError(errorx.ErrCreateUser.Error()).
				WithDebug(err.Error())
		}

		user = &models.User{
			SQLModel: data.SQLModel,
			Email:    data.Email,
			Password: "",
		}
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   user.Id,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := h.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	if err = h.sRepo.SetUserToken(ctx, user.Id, accessToken.GetToken(), subToken, expiredTime); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
