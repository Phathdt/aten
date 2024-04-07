package cmd

import (
	"aten/plugins/dexcomp"
	"aten/plugins/tokenprovider/jwt"
	"fmt"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/component/redisc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phathdt/service-context/component/fiberc"

	"aten/shared/common"

	sctx "github.com/phathdt/service-context"

	"github.com/spf13/cobra"
)

const (
	serviceName = "aten"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(fiberc.New(common.KeyCompFiber)),
		sctx.WithComponent(gormc.NewGormDB(common.KeyCompGorm, "")),
		sctx.WithComponent(jwt.New(common.KeyJwt)),
		sctx.WithComponent(redisc.New(common.KeyCompRedis)),
		sctx.WithComponent(dexcomp.NewDexcomp(common.KeyDex)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 1)

		NewRouter(sc)

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

		// gracefully shutdown
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		_ = sc.Stop()
		logger.Info("Server exited")
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(migrateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
