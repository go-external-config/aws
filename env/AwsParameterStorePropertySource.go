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

const AWSPARAM = "AWSPARAM:"

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
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), AWSPARAM)
		}
	}
	return false
}

func (this *AwsParameterStorePropertySource) Property(key string) string {
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			parameterName := fmt.Sprint(this.environment.ResolveRequiredPlaceholders(source.Property(key)[len(AWSPARAM):]))
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
