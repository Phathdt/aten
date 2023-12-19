package fiberauth

import (
	"aten/plugins/dexcomp"
	"aten/shared/common"
	"errors"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"net/http"
)

type CallbackParams struct {
	Code string `json:"code"`
}

func OauthCallback(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dex := sc.MustGet(common.KeyDex).(dexcomp.DexComponent)

		oauthConfig, err := dex.GetOauthConfig()
		if err != nil {
			panic(err)
		}

		var data CallbackParams

		if err := ctx.QueryParser(&data); err != nil {
			panic(err)
		}
		oauth2Token, err := oauthConfig.Exchange(ctx.Context(), data.Code)
		if err != nil {
			panic(err)
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			panic(errors.New("missing raw id token"))
		}

		idTokenVerifier, err := dex.GetIdTokenProvider()
		if err != nil {
			panic(err)
		}
		// Parse and verify ID Token payload.
		idToken, err := idTokenVerifier.Verify(ctx.Context(), rawIDToken)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(claims))
	}
}
