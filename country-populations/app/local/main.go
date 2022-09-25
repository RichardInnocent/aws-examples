package main

import (
	"country-populations/app/db"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewStaticCredentials("fake", "accessKeyId", "secretKeyId"),
	})
	if err != nil {
		panic(fmt.Errorf("could not create DynamoDB session: %w", err))
	}
	testClient := dynamodb.New(sess)
	testClient.Endpoint = "http://localhost:8000"

	populationStore := db.NewDynamoDBPopulationStore(testClient, "country-populations")
	if err = populationStore.CreateTable(); err != nil {
		panic(fmt.Errorf("could not create country-populations table: %w", err))
	}
	if err = populationStore.PutPopulation("united kingdom", 67_220_000); err != nil {
		panic(fmt.Errorf("could not save population of United Kingdom: %w", err))
	}
	if err = populationStore.PutPopulation("france", 67_390_000); err != nil {
		panic(fmt.Errorf("could not save population of France: %w", err))
	}
	if err = populationStore.PutPopulation("spain", 47_350_000); err != nil {
		panic(fmt.Errorf("could not save population of Spain: %w", err))
	}
	if err = populationStore.PutPopulation("india", 1_380_000_000); err != nil {
		panic(fmt.Errorf("could not save population of India: %w", err))
	}
	if err = populationStore.PutPopulation("china", 1_402_000_000); err != nil {
		panic(fmt.Errorf("could not save population of China: %w", err))
	}
}
