package fiberauth

import (
	"aten/module/handlers"
	"aten/module/models"
	"aten/module/storage"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"aten/shared/validation"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/component/redisc"
	"github.com/phathdt/service-context/core"
	"net/http"
)

func SignUp(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !common.AllowSignup {
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "not allow",
			})
		}

		var p models.SignupRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)
		rdClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()

		sqlStorage := storage.NewSqlStorage(db)
		sessionStore := storage.NewSessionStore(rdClient)
		hdl := handlers.NewSignupHdl(sqlStorage, sessionStore, tokenProvider)

		token, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(token))
	}
}
