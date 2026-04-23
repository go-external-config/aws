package env

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/env"
	"github.com/go-external-config/go/util/optional"
)

const AWSSM = "AWSSM:"

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
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), AWSSM)
		}
	}
	return false
}

func (this *AwsSecretsManagerPropertySource) Property(key string) string {
	for _, source := range this.environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return this.getSecretValue(source.Property(key)[len(AWSSM):])
		}
	}
	panic(err.NewIllegalArgumentException("No value present for " + key))
}

func (this *AwsSecretsManagerPropertySource) getSecretValue(secretName string) string {
	result := optional.OfCommaErr(this.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	})).OrElsePanic("Cannot get AWS secret")
	return *result.SecretString
}

func (this *AwsSecretsManagerPropertySource) newClient() *secretsmanager.Client {
	config := optional.OfCommaErr(config.LoadDefaultConfig(context.Background())).OrElsePanic("Cannot load AWS config")
	return secretsmanager.NewFromConfig(config)
}

func (this *AwsSecretsManagerPropertySource) Properties() map[string]string {
	return nil
}
