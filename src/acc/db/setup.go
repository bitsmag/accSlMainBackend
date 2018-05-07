package db

import (
	"acc/types"
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
)

// DataBase is the driver for database functionality
var DataBase *scribble.Driver

// SetUp creates the database driver and corresponding files
func SetUp(path string) error {
	var err error
	DataBase, err = scribble.New(path, nil)
	if err != nil {
		return err
	}
	// Balance is stored in dynamoDb therefor only to log file needs to be created at startup
	if logsFileExists := canReadLogsFromDb(); !logsFileExists {
		var entries []types.LogEntry // create empty log storage
		if err := forceWriteLogs(entries); err != nil {
			return fmt.Errorf("couldn't set log storage: %v", err)
		}
	}
	return nil
}

func canReadLogsFromDb() bool {
	if _, err := ReadLogs(); err != nil {
		return false
	}
	return true
}

type balanceObj struct {
	AccId   string `json:"AccId"`
	Balance string `json:"Balance"`
}
