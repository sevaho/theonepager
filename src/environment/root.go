package environment

import (
	"log"

	"github.com/rs/zerolog"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	RELEASE                    string        `default:"0.0.1"`
	LOG_LEVEL                  zerolog.Level `default:"1"`
	IS_DEVELOPMENT             bool          `default:"false"`
	STATIC_DIRECTORY           string        `default:"static"`
	TEMPLATES_DIRECTORY        string        `default:"templates"`
	HOT_RELOAD_TEMPLATE_FOLDER bool
	CONFIG_FILE_PATH           string
}

func New() *Environment {

	env := Environment{}

	// parse env vars to struct
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("Failed to decode env vars to struct environment/root.go: %s", err)
	}

	return &env
}
