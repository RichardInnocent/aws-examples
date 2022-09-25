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
localhost:8080?country=spain
```

For the full list of supported countries, see [this file](local/main.go).