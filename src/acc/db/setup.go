package db

import (
	"fmt"

	"github.com/bitsmag/accSlMainBackend/src/acc/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// SetUp creates the required tables in dynamoDb. Caution: overrides existing tables!
func SetUpTables() error {
	fmt.Println("SETTING UP NEW DATABASE TABLES")

	var err error

	balanceDeleteTableInput := getBalanceDeleteTableInput()
	err = deleteTable(&balanceDeleteTableInput)
	if err != nil {
		return fmt.Errorf("couldn't delete table acc_balance: %v", err.Error())
	}
	fmt.Println("*Deleted old balance table")
	balanceCreateTableInput := getBalanceCreateTableInput()
	err = createTable(&balanceCreateTableInput)
	if err != nil {
		return fmt.Errorf("couldn't create table acc_balance: %v", err.Error())
	}
	fmt.Println("*Created new balance table")

	logsDeleteTableInput := getLogsDeleteTableInput()
	err = deleteTable(&logsDeleteTableInput)
	if err != nil {
		return fmt.Errorf("couldn't delete table acc_logs: %v", err.Error())
	}
	fmt.Println("*Deleted old logs table")
	logsCreateTableInput := getLogsCreateTableInput()
	err = createTable(&logsCreateTableInput)
	if err != nil {
		return fmt.Errorf("couldn't create table acc_logs: %v", err.Error())
	}
	fmt.Println("*Created new logs table")

	return nil
}

// Inserts initial/required values into tables
func InsertInitialValues() error {
	var err error

	err = forceWriteBalance(0)
	if err != nil {
		return fmt.Errorf("couldn't write init value '0' to balance table: %v", err.Error())
	}
	return nil
}

type balanceObj struct {
	AccId   string `json:"AccId"`
	Balance string `json:"Balance"`
}

type logObj struct {
	BookingId string         `json:"BookingId"`
	LogEntry  types.LogEntry `json:"LogEntry"`
}

func deleteTable(input *dynamodb.DeleteTableInput) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)

	_, err = svc.DeleteTable(input)

	tableNotFoundException := false
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				tableNotFoundException = true
			default:
				return aerr
			}
		}
		if !tableNotFoundException {
			return err
		}
	}

	dti := dynamodb.DescribeTableInput{TableName: input.TableName}
	err = nil
	err = svc.WaitUntilTableNotExists(&dti)
	return err
}

func createTable(input *dynamodb.CreateTableInput) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)

	_, err = svc.CreateTable(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
		return err
	}

	dti := dynamodb.DescribeTableInput{TableName: input.TableName}
	err = nil
	err = svc.WaitUntilTableExists(&dti)
	return err
}

func getBalanceDeleteTableInput() dynamodb.DeleteTableInput {
	input := dynamodb.DeleteTableInput{
		TableName: aws.String("Acc_balances"),
	}
	return input
}

func getLogsDeleteTableInput() dynamodb.DeleteTableInput {
	input := dynamodb.DeleteTableInput{
		TableName: aws.String("Acc_logs"),
	}
	return input
}

func getBalanceCreateTableInput() dynamodb.CreateTableInput {
	input := dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("AccId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("AccId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("Acc_balances"),
	}
	return input
}

func getLogsCreateTableInput() dynamodb.CreateTableInput {
	input := dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("BookingId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("BookingId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("Acc_logs"),
	}
	return input
}
