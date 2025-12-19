package app

func (app *App) addRoutes(config *Config) {
	// * * * * * * * * * * * * * * * * * *
	// SETUP ROUTES
	// * * * * * * * * * * * * * * * * * *
	app.server.GET("/", handleGetApplications(config))
	app.server.GET("/applications", handleGetApplications(config))

	// healthz
	app.server.GET("/healthz", handleGetHealthz())

	// APIS
	app.server.GET("/v1/api/iconcache", handleGetIconCache())
}
