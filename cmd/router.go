package cmd

import (
	"aten/module/transport/fiberauth"
	"aten/shared/common"
	middleware2 "aten/shared/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/fiberc"
	"github.com/phathdt/service-context/component/fiberc/middleware"
)

func NewRouter(sc sctx.ServiceContext) {
	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})

	app.Use(flogger.New(flogger.Config{
		Format: `{"ip":${ip}, "timestamp":"${time}", "status":${status}, "latency":"${latency}", "method":"${method}", "path":"${path}"}` + "\n",
	}))
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.Recover(sc))

	app.Post("/auth/signup", fiberauth.SignUp(sc))
	app.Post("/auth/login", fiberauth.Login(sc))
	app.Get("/auth/connect", fiberauth.OauthConnect(sc))
	app.Get("/auth/callback", fiberauth.OauthCallback(sc))

	app.Use(middleware2.RequiredAuth(sc))

	app.Get("/auth/me", fiberauth.GetMe(sc))
	app.Get("/auth/valid", fiberauth.CheckValid(sc))

	app.Get("/", ping())

	fiberComp := sc.MustGet(common.KeyCompFiber).(fiberc.FiberComponent)
	fiberComp.SetApp(app)
}

func ping() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(&fiber.Map{
			"msg": "pong",
		})
	}
}
