# AwsSecretsManagerPropertySource

AWS Secrets Manager helps you manage, retrieve, and rotate database credentials, application credentials, OAuth tokens, API keys, and other secrets throughout their lifecycles. Many AWS services store and use secrets in Secrets Manager.

Secrets Manager helps you improve your security posture, because you no longer need hard-coded credentials in application source code. Storing the credentials in Secrets Manager helps avoid possible compromise by anyone who can inspect your application or the components. You replace hard-coded credentials with a runtime call to the Secrets Manager service to retrieve credentials dynamically when you need them.

With Secrets Manager, you can configure an automatic rotation schedule for your secrets. This enables you to replace long-term secrets with short-term ones, significantly reducing the risk of compromise. Since the credentials are no longer stored with the application, rotating credentials no longer requires updating your applications and deploying changes to application clients. ([more](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html))

    export AWS_ACCESS_KEY_ID=...
    export AWS_SECRET_ACCESS_KEY=...
    export AWS_REGION=eu-north-1

cmd/app/main.go

    import (
    	aws "github.com/go-external-config/aws/env"
    	"github.com/go-external-config/go/env"
    )

    var _ = env.Instance().WithPropertySource(aws.NewAwsSecretsManagerPropertySource())

    func main() {
    	fmt.Println("db.pass: " + env.Value[string]("${db.pass}"))
    }

cmd/config/application.yaml

    db:
      pass: AWSSM:db-pass-secret-name
