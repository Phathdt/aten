package fiberauth

import (
	"aten/plugins/dexcomp"
	"aten/shared/common"
	"github.com/gofiber/fiber/v2"
	"github.com/jaevor/go-nanoid"
	sctx "github.com/phathdt/service-context"
	"net/url"
)

func OauthConnect(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dex := sc.MustGet(common.KeyDex).(dexcomp.DexComponent)

		oauthConfig, err := dex.GetOauthConfig()
		if err != nil {
			panic(err)
		}

		canonicID, _ := nanoid.Standard(21)
		state := canonicID()

		authUrl := oauthConfig.AuthCodeURL(state)

		u, err := url.Parse(authUrl)
		if err != nil {
			panic(err)
		}

		u.Path += "/github"

		return ctx.Redirect(u.String())
	}
}
