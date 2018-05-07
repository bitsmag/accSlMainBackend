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

// ProcessTransaction bundles all db-operations necessary for a booking
func ProcessTransaction(amount float64, date types.Date, category types.Category) error {
	if err := bookAmount(amount); err != nil {
		return fmt.Errorf("couldn't write booking to database: %v", err)
	}

	logEntry := types.LogEntry{Amount: amount, Date: date, Category: category}
	if err := logBooking(logEntry); err != nil {
		if err := bookAmount(-amount); err != nil { // roll back booking
			return fmt.Errorf("couldn't write logs to database - state of database might be inconsistent: %v", err)
		}
		return fmt.Errorf("couldn't write booking to database: %v", err)
	}

	return nil
}

// ForceWriteBalance overrides the balance with the passed amount
func forceWriteBalance(balance float64) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)

	balanceObj := balanceObj{AccId: "default", Balance: strconv.FormatFloat(balance, 'f', 6, 64)}
	item, err := dynamodbattribute.MarshalMap(balanceObj)

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("acc_balances"),
	}
	_, err = svc.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

// BookAmount adds/substracts the passed amount to/from the balance
func bookAmount(amount float64) error {
	var balance float64
	var err error
	if balance, err = ReadBalance(); err != nil {
		return fmt.Errorf("couldn't read balance from database: %v", err)
	}
	balance += amount
	if err := forceWriteBalance(balance); err != nil {
		return fmt.Errorf("couldn't write balance to database: %v", err)
	}
	return nil
}

// ForceWriteLogs overrides the logs with the passed logEntry
func forceWriteLogs(entries []types.LogEntry) error {
	DataBase.Delete("logs", "default") // error is thrown when file doesnt exist - ignore error
	if err := DataBase.Write("logs", "default", entries); err != nil {
		return err
	}
	return nil
}

// LogBooking writes a logEntry to the log
func logBooking(entry types.LogEntry) error {
	var entries []types.LogEntry
	var err error
	if entries, err = ReadLogs(); err != nil {
		return fmt.Errorf("couldn't read log-entries from database: %v", err)
	}
	entries = append(entries, entry)
	if err := DataBase.Write("logs", "default", entries); err != nil {
		return err
	}
	return nil

}
