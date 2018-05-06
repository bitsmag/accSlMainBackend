package db

import (
	"acc/types"
)

// ReadBalance returns the current ballance of the account
func ReadBalance() (float64, error) {
	var balance float64
	if err := DataBase.Read("balances", "default", &balance); err != nil {
		return 0, err
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
