package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
	"time"
)

const (
	partitionKeyAttributeName = "partitionKey"
	sortKeyAttributeName      = "sortKey"
	gsi1AttributeName         = "gsi1"
	lsi1AttributeName         = "lsi1"
)

type AddressItem struct {
	userId           string
	yearOfOccupation int
	addressLine1     string
	addressLine2     string
	postcode         string
	country          string
}

func main() {
	tableName := "userAddresses"
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("could not create config: %w", err))
	}
	dynamoDBClient := dynamodb.NewFromConfig(cfg)
	if err := createTableIfNotExists(ctx, tableName, dynamoDBClient); err != nil {
		panic(fmt.Errorf("could not create DynamoDB table: %w", err))
	}
	if err := populateTable(tableName, dynamoDBClient); err != nil {
		panic(fmt.Errorf("could not populate table: %w", err))
	}
}

func createTableIfNotExists(
	ctx context.Context,
	tableName string,
	dynamoDBClient *dynamodb.Client,
) error {
	_, describeErr := dynamoDBClient.DescribeTable(
		ctx,
		&dynamodb.DescribeTableInput{TableName: addressOfString(tableName)},
	)
	if describeErr == nil {
		fmt.Println("Table already exists. Skipping...")
		return nil
	}
	fmt.Printf("Could not describe the table: %v\nAttempting to create table...\n", describeErr)
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: addressOfString(partitionKeyAttributeName), // partition key
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: addressOfString(sortKeyAttributeName), // sort key
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: addressOfString(lsi1AttributeName), // LSI
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: addressOfString(gsi1AttributeName), // GSI
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: addressOfString(partitionKeyAttributeName),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: addressOfString(sortKeyAttributeName),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: addressOfString("gsi1"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: addressOfString(gsi1AttributeName),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeKeysOnly,
				},
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: addressOfString("lsi1"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: addressOfString(partitionKeyAttributeName),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: addressOfString(lsi1AttributeName),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeKeysOnly,
				},
			},
		},
		TableClass: types.TableClassStandard,
		TableName:  addressOfString(tableName),
	}
	_, err := dynamoDBClient.CreateTable(ctx, createTableInput)

	if err != nil {
		fmt.Println(
			"Sleeping for 5 seconds to give the table some time to create before inserting items",
		)
		time.Sleep(5 * time.Second)
	}
	return err
}

func populateTable(tableName string, dynamoDBService *dynamodb.Client) error {
	items := []AddressItem{
		{
			userId:           "1",
			yearOfOccupation: 2012,
			addressLine1:     "1027 Tribecca Road",
			addressLine2:     "Hirsten",
			postcode:         "JD71CG",
			country:          "United Kingdom",
		},
		{
			userId:           "1",
			yearOfOccupation: 2016,
			addressLine1:     "76 Wyjo Avenue",
			addressLine2:     "Pollena",
			postcode:         "LB47RF",
			country:          "Spain",
		},
		{
			userId:           "1",
			yearOfOccupation: 2020,
			addressLine1:     "1 Heather Road",
			addressLine2:     "Lannert",
			postcode:         "RT659PL",
			country:          "United Kingdom",
		},
		{
			userId:           "1",
			yearOfOccupation: 2021,
			addressLine1:     "115 Scarsdale Terrace",
			addressLine2:     "Crannage",
			postcode:         "BN761FG",
			country:          "United Kingdom",
		},
		{
			userId:           "2",
			yearOfOccupation: 2014,
			addressLine1:     "19 Cookie Close",
			addressLine2:     "Forten",
			postcode:         "AS129TY",
			country:          "Italy",
		},
		{
			userId:           "2",
			yearOfOccupation: 2019,
			addressLine1:     "64 Valencia Place",
			addressLine2:     "Minnert",
			postcode:         "ES94NI07",
			country:          "Spain",
		},
		{
			userId:           "3",
			yearOfOccupation: 2015,
			addressLine1:     "115 Scarsdale Terrace",
			addressLine2:     "Crannage",
			postcode:         "BN761FG",
			country:          "United Kingdom",
		},
		{
			userId:           "3",
			yearOfOccupation: 2017,
			addressLine1:     "2 Beach Avenue",
			addressLine2:     "Trusset",
			postcode:         "BA142RT",
			country:          "United Kingdom",
		},
		{
			userId:           "3",
			yearOfOccupation: 2019,
			addressLine1:     "Chan Cottage",
			addressLine2:     "West Covena",
			postcode:         "RO329ZX",
			country:          "United Kingdom",
		},
		{
			userId:           "3",
			yearOfOccupation: 2022,
			addressLine1:     "6 Nathaniel Court",
			addressLine2:     "Solitun",
			postcode:         "PL14ER",
			country:          "United Kingdom",
		},
		{
			userId:           "4",
			yearOfOccupation: 2013,
			addressLine1:     "10 Whitefeather Crescent",
			addressLine2:     "Brandish",
			postcode:         "DA74TP",
			country:          "United Kingdom",
		},
	}

	for _, item := range items {
		if err := insertItem(tableName, dynamoDBService, item); err != nil {
			return fmt.Errorf("could not insert item for user ID %s: %w", item.userId, err)
		}
	}
	return nil
}

func insertItem(tableName string, dynamoDBService *dynamodb.Client, item AddressItem) error {
	_, err := dynamoDBService.PutItem(context.Background(),
		&dynamodb.PutItemInput{
			TableName: addressOfString(tableName),
			Item: map[string]types.AttributeValue{
				partitionKeyAttributeName: &types.AttributeValueMemberS{
					Value: "userId/" + item.userId,
				},
				sortKeyAttributeName: &types.AttributeValueMemberS{
					Value: "year/" + strconv.Itoa(item.yearOfOccupation),
				},
				gsi1AttributeName: &types.AttributeValueMemberS{
					Value: "postcode/" + item.postcode,
				},
				lsi1AttributeName: &types.AttributeValueMemberS{
					Value: "country/" + item.country,
				},
				"userId": &types.AttributeValueMemberS{
					Value: item.userId,
				},
				"yearOfOccupation": &types.AttributeValueMemberS{
					Value: strconv.Itoa(item.yearOfOccupation),
				},
				"addressLine1": &types.AttributeValueMemberS{
					Value: item.addressLine1,
				},
				"addressLine2": &types.AttributeValueMemberS{
					Value: item.addressLine2,
				},
				"postcode": &types.AttributeValueMemberS{
					Value: item.postcode,
				},
				"country": &types.AttributeValueMemberS{
					Value: item.country,
				},
			},
			ReturnValues: types.ReturnValueNone,
		},
	)
	return err
}

func addressOfString(value string) *string {
	return &value
}
