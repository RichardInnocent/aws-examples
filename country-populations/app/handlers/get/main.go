package main

import (
	"country-populations/app/db"
	"country-populations/app/handlers"
	"fmt"
	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joerdav/zapray"
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
	handler := handlers.NewIndexHandler(logger, populationStore)
	algnhsa.ListenAndServe(handler, nil)
}

func getDynamoDBClient(sess *session.Session) *dynamodb.DynamoDB {
	client := dynamodb.New(sess)
	return client
}

func getSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, fmt.Errorf("could not create session: %w", err)
	}
	return sess, nil
}
