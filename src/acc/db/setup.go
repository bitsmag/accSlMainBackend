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
	if balanceFileExists := canReadBalanceFromDb(); !balanceFileExists {
		if err := forceWriteBalance(0); err != nil { // create "0" balance storage
			return fmt.Errorf("couldn't set balance storage: %v", err)
		}
		var entries []types.LogEntry // create empty log storage
		if err := forceWriteLogs(entries); err != nil {
			return fmt.Errorf("couldn't set log storage: %v", err)
		}
	}
	return nil
}

func canReadBalanceFromDb() bool {
	if _, err := ReadBalance(); err != nil {
		return false
	}
	return true
}
