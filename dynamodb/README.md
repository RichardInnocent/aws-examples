# DynamoDB
This is a lab on DynamoDB.

## Setup
For the setup, configure the credentials in `00-set-credentials.sh` and run the command.

Then, if the example table hasn't already been created, run
```shell
go run main.go
```

This will create the table and populate it with a few items.

## Example queries

### Getting by partition key and sort key exact match
```shell
aws dynamodb get-item \
  --table-name userAddresses \
  --key '{"partitionKey": {"S": "userId/1"}, "sortKey": {"S": "year/2021"}}'
```

### Getting by partition key only
```shell
aws dynamodb query \
  --table-name userAddresses \
  --key-condition-expression "partitionKey = :pk" \
  --expression-attribute-values '{":pk": {"S": "userId/1"}}'
```

### Getting by partition key and sort key range
```shell
aws dynamodb query \
  --table-name userAddresses \
  --key-condition-expression "partitionKey = :pk AND sortKey BETWEEN :start AND :end" \
  --expression-attribute-values '{":pk": {"S": "userId/1"}, ":start": {"S": "year/2015"}, ":end": {"S": "year/2021"}}'
```

### Getting by partition key and sort key range
```shell
aws dynamodb query \
  --table-name userAddresses \
  --index-name lsi1 \
  --key-condition-expression "partitionKey = :pk AND lsi1 = :lsi" \
  --expression-attribute-values '{":pk": {"S": "userId/1"}, ":lsi": {"S": "country/Spain"}}'
```

### Searching by global secondary index
```shell
aws dynamodb query \
  --table-name userAddresses \
  --index-name gsi1 \
  --key-condition-expression "gsi1 = :gsi1" \
  --expression-attribute-values '{":gsi1": {"S": "postcode/BN761FG"}}'
```
