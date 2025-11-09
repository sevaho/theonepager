package main

import (
	"context"

	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	app "github.com/sevaho/theonepager/src"
	"github.com/sevaho/theonepager/src/environment"
	"github.com/sevaho/theonepager/src/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	// parse arguments (flags)
	var (
		serve          = pflag.Bool("serve", false, "Serve the application.")
		port           = pflag.IntP("port", "p", 3000, "Which port to run on.")
		configfilepath = pflag.StringP("config", "c", "config.yaml", "Path to config file.")
	)

	pflag.Parse()

	if *serve {
		env := environment.New()

		if configfilepath != nil {
			env.CONFIG_FILE_PATH = *configfilepath
		}

		ctx := context.Background()
		app.Run(ctx, *port, env)
	} else {
		logger.Logger.Warn().Msg("No flags given, exiting!")
		pflag.PrintDefaults()
	}
}
