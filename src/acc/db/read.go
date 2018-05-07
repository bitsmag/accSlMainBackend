package db

import (
	"acc/types"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ReadBalance returns the current ballance of the account
func ReadBalance() (float64, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("acc_balances"),
		Key: map[string]*dynamodb.AttributeValue{
			"AccId": {S: aws.String("default")},
		},
	})
	if err != nil {
		return 0, err
	}

	balanceObj := balanceObj{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &balanceObj)
	balance, err := strconv.ParseFloat(balanceObj.Balance, 64)

	if err != nil {
		return 0, fmt.Errorf("Failed to unmarshal and parse balance: %v", err)
	}

	return balance, nil
}

// ReadLogs returns all logEntries
func ReadLogs() ([]types.LogEntry, error) {
	var entries []types.LogEntry
	if err := DataBase.Read("logs", "default", &entries); err != nil {
		return entries, err
	}
	return entries, nil
}
