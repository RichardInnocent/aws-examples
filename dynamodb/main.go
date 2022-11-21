package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	tableName := "userAddresses"

	dynamoDBService := dynamodb.New(sess)
	if err := createTable(tableName, dynamoDBService); err != nil {
		panic(fmt.Errorf("could not create DynamoDB table: %w", err))
	}
	if err := populateTable(tableName, dynamoDBService); err != nil {
		panic(fmt.Errorf("could not populate table: %w", err))
	}
}

func createTable(tableName string, dynamoDBService *dynamodb.DynamoDB) error {
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: addressOfString(partitionKeyAttributeName), // partition key
				AttributeType: addressOfString("S"),
			},
			{
				AttributeName: addressOfString(sortKeyAttributeName), // sort key
				AttributeType: addressOfString("S"),
			},
			{
				AttributeName: addressOfString(lsi1AttributeName), // LSI
				AttributeType: addressOfString("S"),
			},
			{
				AttributeName: addressOfString(gsi1AttributeName), // GSI
				AttributeType: addressOfString("S"),
			},
		},
		BillingMode: addressOfString("PAY_PER_REQUEST"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: addressOfString(partitionKeyAttributeName),
				KeyType:       addressOfString("HASH"),
			},
			{
				AttributeName: addressOfString(sortKeyAttributeName),
				KeyType:       addressOfString("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: addressOfString("postcode"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: addressOfString(gsi1AttributeName),
						KeyType:       addressOfString("S"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: addressOfString(dynamodb.ProjectionTypeKeysOnly),
				},
			},
		},
		LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			{
				IndexName: addressOfString("country"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: addressOfString(lsi1AttributeName),
						KeyType:       addressOfString("S"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: addressOfString(dynamodb.ProjectionTypeKeysOnly),
				},
			},
		},
		TableClass: addressOfString("STANDARD"),
		TableName:  addressOfString(tableName),
	}
	_, err := dynamoDBService.CreateTable(createTableInput)
	return err
}

func populateTable(tableName string, dynamoDBService *dynamodb.DynamoDB) error {
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

func insertItem(tableName string, dynamoDBService *dynamodb.DynamoDB, item AddressItem) error {
	_, err := dynamoDBService.PutItem(&dynamodb.PutItemInput{
		ExpressionAttributeNames:  nil,
		ExpressionAttributeValues: nil,
		Item: map[string]*dynamodb.AttributeValue{
			partitionKeyAttributeName: {S: addressOfString("userId/" + item.userId)},
			sortKeyAttributeName:      {S: addressOfString("year/" + strconv.Itoa(item.yearOfOccupation))},
			gsi1AttributeName:         {S: addressOfString("postcode/" + item.postcode)},
			lsi1AttributeName:         {S: addressOfString("country/" + item.country)},
			"userId":                  {S: addressOfString(item.userId)},
			"yearOfOccupation":        {N: addressOfString(strconv.Itoa(item.yearOfOccupation))},
			"addressLine1":            {S: addressOfString(item.addressLine1)},
			"addressLine2":            {S: addressOfString(item.addressLine2)},
			"postcode":                {S: addressOfString(item.postcode)},
			"country":                 {S: addressOfString(item.country)},
		},
		ReturnValues: nil,
		TableName:    addressOfString(tableName),
	})
	return err
}

func addressOfString(value string) *string {
	return &value
}
