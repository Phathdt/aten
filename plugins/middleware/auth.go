package middleware

import (
	sctx "github.com/CaliberVB/service-context"
	"github.com/CaliberVB/service-context/core"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"strings"
	common "telifi-common"
	"telifi-plugins/tokenprovider"
)

func extractTokenFromHeaderString(headers []string) (string, error) {
	if len(headers) == 0 {
		return "", errors.New("missing token")
	}
	//"Authorization" : "Bearer {token}"

	parts := strings.Split(headers[0], " ")

	if len(parts) == 0 {
		return "", errors.New("missing token")
	}

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("wrong authen header")
	}

	return parts[1], nil
}

func RequiredAuth(sc sctx.ServiceContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		token, err := extractTokenFromHeaderString(headers["Authorization"])

		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		c.Context().SetUserValue("userId", payload.GetUserId())
		return c.Next()
	}
}
