package app

func (app *App) addRoutes(config *Config) {
	// * * * * * * * * * * * * * * * * * *
	// SETUP ROUTES
	// * * * * * * * * * * * * * * * * * *
	app.server.GET("/", handleGetIndex(config))

	// healthz
	app.server.GET("/healthz", handleGetIndex(config))

	// APIS
	app.server.GET("/v1/api/iconcache", handleGetIconCache())
}
