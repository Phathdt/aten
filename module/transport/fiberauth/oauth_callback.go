package fiberauth

import (
	"aten/module/handlers"
	"aten/module/storage"
	"aten/plugins/dexcomp"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"fmt"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/component/redisc"
	"github.com/phathdt/service-context/core"
	"net/http"
)

type CallbackParams struct {
	Code string `json:"code"`
}

func OauthCallback(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p CallbackParams

		logger := sctx.GlobalLogger().GetLogger("service")

		if err := ctx.QueryParser(&p); err != nil {
			panic(err)
		}
		dex := sc.MustGet(common.KeyDex).(dexcomp.DexComponent)
		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)
		rdClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()

		sqlStorage := storage.NewSqlStorage(db)
		sessionStore := storage.NewSessionStore(rdClient)
		hdl := handlers.NewOauthCallbackHdl(sqlStorage, sessionStore, dex, tokenProvider)
		res, err := hdl.Response(ctx.Context(), p.Code)
		if err != nil {
			logger.Error(err)

			if dex.GetRedirect() {
				return ctx.Redirect(fmt.Sprintf("%s", dex.GetClientErrEndpoint()), http.StatusMovedPermanently)
			}

			return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(err))
		}

		if dex.GetRedirect() {
			return ctx.Redirect(fmt.Sprintf("%s?token=%s", dex.GetClientEndpoint(), res.GetToken()), http.StatusMovedPermanently)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(res))
	}
}
