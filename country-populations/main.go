package main

import (
	"country-populations/db"
	"country-populations/handlers"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joerdav/zapray"
	"net/http"
	"os"
)

func main() {
	sess, err := getSession("eu-west-1")
	if err != nil {
		panic(fmt.Errorf("could not start server: %w", err))
	}
	dynamoDBClient := getDynamoDBClient(sess)
	populationStore := db.NewDynamoDBPopulationStore(dynamoDBClient, "country-populations")
	logger, err := zapray.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("could not create logger: %w", err))
	}
	http.Handle("/", handlers.NewIndexHandler(logger, populationStore))
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(fmt.Errorf("could not start HTTP server: %w", err))
	}
}

func getDynamoDBClient(sess *session.Session) *dynamodb.DynamoDB {
	client := dynamodb.New(sess)
	if isLocal() {
		client.Endpoint = "http://localhost:8000"
	}
	return client
}

func getSession(region string) (*session.Session, error) {
	if isLocal() {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("eu-west-1"),
			Credentials: credentials.NewStaticCredentials("fake", "accessKeyId", "secretKeyId"),
		})
		if err != nil {
			return nil, fmt.Errorf("could not create local session: %w", err)
		}
		return sess, nil
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, fmt.Errorf("could not create session: %w", err)
	}
	return sess, nil
}

func isLocal() bool {
	for _, arg := range os.Args {
		if arg == "--local" {
			return true
		}
	}
	return false
}
