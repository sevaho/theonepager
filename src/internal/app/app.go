package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sevaho/theonepager/src/environment"
	"github.com/sevaho/theonepager/src/pkg/logger"
	"github.com/sevaho/theonepager/src/pkg/renderer"
	"github.com/sevaho/theonepager/src/web"
	"github.com/unrolled/secure"

	"github.com/sevaho/livereload"
)

type App struct {
	port     int
	server   *echo.Echo
	renderer *renderer.RenderEngine
	env      *environment.Environment
	config   *Config
}

// TODO:  <03-05-25, Sebastiaan Van Hoecke> // This should return a pointer to echo
func New(env *environment.Environment) *App {
	server := echo.New()
	server.Debug = env.IS_DEVELOPMENT

	if env.IS_DEVELOPMENT {
		server.StaticFS("/static", os.DirFS(env.STATIC_DIRECTORY))
	} else {
		staticFS, _ := fs.Sub(web.Static, "static")
		server.StaticFS("/static", staticFS)
	}
	server.HideBanner = true
	server.HidePort = true

	// * * * * * * * * * * * * * * * * * *
	// CONFIG
	// * * * * * * * * * * * * * * * * * *
	config, err := parseConfig(env.CONFIG_FILE_PATH)
	if err != nil {
		panic(err)
	}

	// * * * * * * * * * * * * * * * * * *
	// RENDER ENGINE
	// * * * * * * * * * * * * * * * * * *
	renderer := renderer.NewRenderEngine(env.IS_DEVELOPMENT, env.TEMPLATES_DIRECTORY, env.RELEASE, web.Templates)
	server.Renderer = renderer

	// * * * * * * * * * * * * * * * * * *
	// DEPENDENCIES
	// * * * * * * * * * * * * * * * * * *

	app := App{
		server:   server,
		renderer: renderer,
		env:      env,
		config:   config,
	}

	// * * * * * * * * * * * * * * * * * *
	// MIDDLEWARE
	// * * * * * * * * * * * * * * * * * *
	server.Use(middleware.RequestID())
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:     true,
		IsDevelopment: env.IS_DEVELOPMENT,
	})
	server.Use(echo.WrapMiddleware(secureMiddleware.Handler))

	if env.IS_DEVELOPMENT {
		server.Use(livereload.LiveReload(server, logger.Logger, env.TEMPLATES_DIRECTORY, env.STATIC_DIRECTORY))
	}

	// * * * * * * * * * * * * * * * * * *
	// ROUTES
	// * * * * * * * * * * * * * * * * * *
	app.addRoutes(config)

	return &app
}

func (app *App) Serve(port int) {
	go func() {
		logger.Logger.Info().Msgf("[HTTP SERVER] Running on http://localhost:%d", port)
		err := app.server.Start(fmt.Sprint(":", port))
		if err != nil {
			logger.Logger.Warn().Err(err).Stack().Msgf("[HTTP SERVER] closed unexpectedly, reason: %s", err)
		}
	}()
}

func (app *App) ShutDown(ctx context.Context) {
	err := app.server.Shutdown(ctx)
	if err != nil {
		logger.Logger.Error().Err(err).Stack().Msgf("[HTTP SERVER] Error shutting down the server: %s", err)
	}
	logger.Logger.Info().Msgf("[HTTP SERVER] shut down.")
}
