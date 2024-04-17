package root

import (
	"fmt"
	"simple-api-go/internal/app/context"
	"simple-api-go/variable"
)

const (
	FlagEnv      string = "env"
	FlagEnvShort string = "e"
)

var (
	App         context.IAppContext
	ConfigPaths = [...]string{
		"./*.toml",
		"./params/*.toml",
		"/opt/" + variable.AppName + "/params/*.toml",
	}

	// FlagEnvValue default environment in not defined in args flag
	FlagEnvValue = EnvLocal
)

type AppEnv string

const (
	AllowedEnvInfo        = `"local", "development", "staging", "uat" or "production"`
	EnvLocal       AppEnv = "local"
	EnvDevelopment AppEnv = "development"
	EnvStaging     AppEnv = "staging"
	EnvUat         AppEnv = "uat"
	EnvProduction  AppEnv = "production"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *AppEnv) String() string {
	return string(*e)
}

// Set must have pointer receiver, so it doesn't change the value of a copy
func (e *AppEnv) Set(v string) error {
	allowed := map[AppEnv]struct{}{
		EnvLocal:       {},
		EnvDevelopment: {},
		EnvStaging:     {},
		EnvUat:         {},
		EnvProduction:  {},
	}

	if _, ok := allowed[AppEnv(v)]; !ok {
		return fmt.Errorf(`must be one of %s`, AllowedEnvInfo)
	}

	*e = AppEnv(v)
	return nil
}

// Type is only used in help text
func (e *AppEnv) Type() string {
	return "AppEnv"
}
