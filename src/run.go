package src

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	"github.com/sevaho/theonepager/src/environment"
	"github.com/sevaho/theonepager/src/internal/app"
	"github.com/sevaho/theonepager/src/pkg/logger"
)

func Run(ctx context.Context, port int, env *environment.Environment) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Setup Logger
	logger.Init(env.IS_DEVELOPMENT, env.LOG_LEVEL)

	application := app.New(env)

	application.Serve(port)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		time.Sleep(1)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		application.ShutDown(shutdownCtx)
	}()
	wg.Wait()
	return nil
}
