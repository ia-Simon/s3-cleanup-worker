package core

import (
	"github.com/ilyakaznacheev/cleanenv"
)

var Settings struct {
	ENVIRONMENT           string `env:"ENVIRONMENT"`
	DATABASE_URL          string `env:"DATABASE_URL"`
	AWS_ACCESS_KEY_ID     string `env:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `env:"AWS_SECRET_ACCESS_KEY"`
	AWS_DEFAULT_REGION    string `env:"AWS_DEFAULT_REGION"`
	AWS_FILES_BUCKET      string `env:"AWS_FILES_BUCKET"`
}

func init() {
	err := cleanenv.ReadEnv(&Settings)
	if err != nil {
		panic(err)
	}

	initAws()
	initDb()
}
