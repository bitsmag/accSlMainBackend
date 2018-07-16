package db

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bitsmag/accSlMainBackend/src/acc/types"

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
	tablenameBalances := os.Getenv("TABLENAME_BALANCES")
	if len(tablenameBalances) == 0 {
		tablenameBalances = "Acc_balances"
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)

	balanceObj := balanceObj{AccId: "default", Balance: strconv.FormatFloat(balance, 'f', 6, 64)}
	item, err := dynamodbattribute.MarshalMap(balanceObj)

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tablenameBalances),
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

// logBooking writes a logEntry to the log
func logBooking(entry types.LogEntry) error {
	tablenameLog := os.Getenv("TABLENAME_LOG")
	if len(tablenameLog) == 0 {
		tablenameLog = "Acc_logs"
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := dynamodb.New(sess)

	bookingId := randomString(14)
	logObj := logObj{BookingId: bookingId, LogEntry: entry}
	item, err := dynamodbattribute.MarshalMap(logObj)

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tablenameLog),
	}
	_, err = svc.PutItem(input)

	if err != nil {
		return err
	}
	return nil
}

// https://www.calhoun.io/creating-random-strings-in-go/
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
