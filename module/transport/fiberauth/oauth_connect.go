package fiberauth

import (
	"aten/module/handlers"
	"aten/plugins/dexcomp"
	"aten/shared/common"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
)

type oauthConnectParams struct {
	ConnectorId string `query:"connector_id"`
}

func OauthConnect(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p oauthConnectParams

		if err := ctx.QueryParser(&p); err != nil {
			panic(err)
		}

		dex := sc.MustGet(common.KeyDex).(dexcomp.DexComponent)

		hdl := handlers.NewOauthConnectHdl(dex)
		url, err := hdl.Response(ctx.Context(), p.ConnectorId)
		if err != nil {
			panic(err)
		}

		return ctx.Redirect(url)
	}
}
