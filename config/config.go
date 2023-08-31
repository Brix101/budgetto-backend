package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	DATABASE_URL string `validate:"required" mapstructure:"DATABASE_URL"`
	TOKEN_SECRET string `validate:"required" mapstructure:"TOKEN_SECRET"`
	SALT         string `validate:"required" mapstructure:"SALT"`
	PORT         string `mapstructure:"PORT"`
}

var config EnvVars

func GetConfig() *EnvVars {
	return &config
}

func init() {
	env := os.Getenv("GO_ENV")

	if env == "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	envType := reflect.TypeOf(&config).Elem()
	envValue := reflect.ValueOf(&config).Elem()
	for i := 0; i < envType.NumField(); i++ {
		field := envType.Field(i)
		fieldValue := envValue.Field(i)

		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			envVarName := field.Tag.Get("mapstructure")
			envVarValue := os.Getenv(envVarName)
			if envVarValue != "" {
				fieldValue.SetString(envVarValue)
			}
		}
	}

	if config.PORT == "" {
		config.PORT = "5000"
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		errs := err.(validator.ValidationErrors)
		for i, e := range errs {
			if i == len(errs)-1 {
				log.Fatalf("❌ Invalid environment variables: %v\n", e.Field())
			} else {
				log.Printf("❌ Invalid environment variables: %v\n", e.Field())
			}
		}
	}
}
