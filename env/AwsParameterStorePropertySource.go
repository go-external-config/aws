package env

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/env"
	"github.com/go-external-config/go/util/optional"
)

const AWSPARAM_KEY_PREFIX = "AWSPARAM."
const AWSPARAM_VALUE_PREFIX = "AWSPARAM:"

type AwsParameterStorePropertySource struct {
	environment *env.Environment
	client      *ssm.Client
}

func NewAwsParameterStorePropertySource() *AwsParameterStorePropertySource {
	ps := &AwsParameterStorePropertySource{}
	ps.environment = env.Instance()
	ps.client = ps.newClient()
	return ps
}

func (this *AwsParameterStorePropertySource) Name() string {
	return "AwsParameterStorePropertySource"
}

func (this *AwsParameterStorePropertySource) HasProperty(key string) bool {
	if strings.HasPrefix(key, AWSPARAM_KEY_PREFIX) {
		return true
	}
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), AWSPARAM_VALUE_PREFIX)
		}
	}
	return false
}

func (this *AwsParameterStorePropertySource) Property(key string) string {
	if strings.HasPrefix(key, AWSPARAM_KEY_PREFIX) {
		parameterName := fmt.Sprint(this.environment.ResolveRequiredPlaceholders(key[len(AWSPARAM_KEY_PREFIX):]))
		return this.getParameterValue(parameterName)
	}
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			parameterName := fmt.Sprint(this.environment.ResolveRequiredPlaceholders(source.Property(key)[len(AWSPARAM_VALUE_PREFIX):]))
			return this.getParameterValue(parameterName)
		}
	}
	panic(err.NewIllegalArgumentException("No value present for " + key))
}

func (this *AwsParameterStorePropertySource) getParameterValue(parameterName string) string {
	result := optional.OfCommaErr(this.client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true),
	})).OrElsePanic("Cannot get AWS parameter " + parameterName)
	return *result.Parameter.Value
}

func (this *AwsParameterStorePropertySource) newClient() *ssm.Client {
	config := optional.OfCommaErr(config.LoadDefaultConfig(context.Background())).OrElsePanic("Cannot load AWS config")
	return ssm.NewFromConfig(config)
}

func (this *AwsParameterStorePropertySource) Properties() map[string]string {
	return nil
}
