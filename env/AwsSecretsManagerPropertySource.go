package env

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/env"
	"github.com/go-external-config/go/util/optional"
)

const AWSSECRET_KEY_PREFIX = "AWSSECRET."
const AWSSECRET_VALUE_PREFIX = "AWSSECRET:"

type AwsSecretsManagerPropertySource struct {
	environment *env.Environment
	client      *secretsmanager.Client
}

func NewAwsSecretsManagerPropertySource() *AwsSecretsManagerPropertySource {
	ps := &AwsSecretsManagerPropertySource{}
	ps.environment = env.Instance()
	ps.client = ps.newClient()
	return ps
}

func (this *AwsSecretsManagerPropertySource) Name() string {
	return "AwsSecretsManagerPropertySource"
}

func (this *AwsSecretsManagerPropertySource) HasProperty(key string) bool {
	if strings.HasPrefix(key, AWSSECRET_KEY_PREFIX) {
		return true
	}
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), AWSSECRET_VALUE_PREFIX)
		}
	}
	return false
}

func (this *AwsSecretsManagerPropertySource) Property(key string) string {
	if strings.HasPrefix(key, AWSSECRET_KEY_PREFIX) {
		parameterName := fmt.Sprint(this.environment.ResolveRequiredPlaceholders(key[len(AWSSECRET_KEY_PREFIX):]))
		return this.getSecretValue(parameterName)
	}
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			secretName := fmt.Sprint(this.environment.ResolveRequiredPlaceholders(source.Property(key)[len(AWSSECRET_VALUE_PREFIX):]))
			return this.getSecretValue(secretName)
		}
	}
	panic(err.NewIllegalArgumentException("No value present for " + key))
}

func (this *AwsSecretsManagerPropertySource) getSecretValue(secretName string) string {
	result := optional.OfCommaErr(this.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	})).OrElsePanic("Cannot get AWS secret " + secretName)
	if result.SecretString != nil {
		return *result.SecretString
	}
	if result.SecretBinary != nil {
		return base64.StdEncoding.EncodeToString(result.SecretBinary)
	}
	panic(err.NewIllegalStateException(fmt.Sprintf("AWS secret '%s' has neither SecretString nor SecretBinary", secretName)))
}

func (this *AwsSecretsManagerPropertySource) newClient() *secretsmanager.Client {
	config := optional.OfCommaErr(config.LoadDefaultConfig(context.Background())).OrElsePanic("Cannot load AWS config")
	return secretsmanager.NewFromConfig(config)
}

func (this *AwsSecretsManagerPropertySource) Properties() map[string]string {
	return nil
}
