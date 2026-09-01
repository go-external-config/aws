package aws

import (
	"github.com/go-external-config/aws/env"
	config "github.com/go-external-config/go/env"
)

func init() {
	config.RegisterPropertySource(env.NewAwsParameterStorePropertySource())
	config.RegisterPropertySource(env.NewAwsSecretsManagerPropertySource())
}
