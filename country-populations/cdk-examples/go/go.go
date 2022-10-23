package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	// "github.com/aws/jsii-runtime-go"
)

type GoStackProps struct {
	awscdk.StackProps
}

func NewGoStack(scope constructs.Construct, id string, props *GoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	bundlingOptions := &awslambdago.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w"`)},
	}

	getPopulationByCountryFunction := awslambdago.NewGoFunction(
		stack,
		jsii.String("getPopulationByCountry"),
		&awslambdago.GoFunctionProps{
			FunctionName: jsii.String("get-population-by-country"),
			Runtime:      awslambda.Runtime_GO_1_X(),
			Entry:        jsii.String("../../app/handlers/get"),
			Bundling:     bundlingOptions,
			Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
			MemorySize:   jsii.Number(1024),
		},
	)

	countryPopulationsTable := awsdynamodb.NewTable(
		scope,
		jsii.String("country-populations-table"),
		&awsdynamodb.TableProps{
			TableName: jsii.String("country-populations"),
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("country"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
	)

	countryPopulationsTable.GrantReadData(getPopulationByCountryFunction)

	apiGateway := awsapigateway.NewLambdaRestApi(
		scope,
		jsii.String("api-gateway"),
		&awsapigateway.LambdaRestApiProps{
			Handler: getPopulationByCountryFunction,
		},
	)
	countriesGetIntegration := awsapigateway.NewLambdaIntegration(
		getPopulationByCountryFunction,
		&awsapigateway.LambdaIntegrationOptions{},
	)
	apiGateway.Root().AddMethod(
		jsii.String("GET"),
		countriesGetIntegration,
		&awsapigateway.MethodOptions{},
	)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoStack(app, "GoStack", &GoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
