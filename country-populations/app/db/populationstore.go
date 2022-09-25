package db

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

type DynamoDBPopulationStore struct {
	dynamoDBClient *dynamodb.DynamoDB
	tableName      string
}

func NewDynamoDBPopulationStore(
	dynamoDBClient *dynamodb.DynamoDB,
	tableName string,
) DynamoDBPopulationStore {
	return DynamoDBPopulationStore{
		dynamoDBClient: dynamoDBClient,
		tableName:      tableName,
	}
}

type CountryPopulation struct {
	Country    string `json:"country"`
	Population int64  `json:"population"`
}

func (store DynamoDBPopulationStore) GetPopulation(country string) (int64, bool, error) {
	result, err := store.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"country": {
				S: &country,
			},
		},
		TableName: &store.tableName,
	})
	if err != nil {
		return 0, false, fmt.Errorf("could not get item from population database: %w", err)
	}
	if result == nil {
		return 0, false, nil
	}
	populationAttribute := result.Item["population"]
	if populationAttribute == nil {
		return 0, false, nil
	}
	populationAsText := populationAttribute.N
	if populationAsText == nil {
		return 0, false, errors.New("population not present on database entry")
	}
	population, err := strconv.ParseInt(*populationAsText, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("could not convert population to a number: %w", err)
	}
	return population, true, nil
}

func (store DynamoDBPopulationStore) PutPopulation(country string, population int64) error {
	_, err := store.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"country": {
				S: &country,
			},
			"population": {
				N: aws.String(strconv.FormatInt(population, 10)),
			},
		},
		TableName: &store.tableName,
	})
	if err != nil {
		return fmt.Errorf("could not put item: %w", err)
	}
	return nil
}

func (store DynamoDBPopulationStore) CreateTable() error {
	_, err := store.dynamoDBClient.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("country"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("country"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   &store.tableName,
	})
	if err != nil {
		return fmt.Errorf("could not create table: %w", err)
	}
	return nil
}
