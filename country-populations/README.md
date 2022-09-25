# Country Populations
This is a demo application that features reading from a DynamoDB table.

## Running locally
### Run DynamoDB locally
Run the DynamoDB shell in a terminal:
```shell
./run-dynamodb-local.sh
```

### Create and populate the table
```shell
cd local
go run main.go
```

### Run the application
From the root of the project, run:
```shell
go run main.go --local
```

### Accessing the application
The application can be accessed by visiting http://localhost:8080

To switch between countries, set the country in the query parameters, e.g.:
```
http://localhost:8080?country=spain
```

For the full list of supported countries, see [this file](app/local/main.go).

## Deploying as a Lambda function
### Set the credentials
Grab the credentials for the account, put them in [oo-set-credentials.sh](00-set-credentials.sh)
and then run it.

### Bootstrap the environment
The environment needs to be bootstrapped so that CDK CF deployments can run. You only need to do
this once per account.

```shell
cdk bootstrap
```

### Run the CDK deployment
```shell
cd cdk
cdk deploy
```