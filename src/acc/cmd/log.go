package cmd

import (
	"fmt"

	"github.com/bitsmag/accSlMainBackend/src/acc/db"
	"github.com/bitsmag/accSlMainBackend/src/acc/types"
	"github.com/bitsmag/accSlMainBackend/src/acc/useroutput"
)

func LogCmdHandler(logOrder types.Order, logDate types.Date, logCategory types.Category) (string, error) {
	var entries []types.LogEntry
	var err error
	if entries, err = db.ReadLogs(); err != nil {
		return "", fmt.Errorf("couldn't read log-entries from database: %v", err)
	}

	filteredEntries := filterEntries(entries, logDate, logCategory)

	logOutput := useroutput.CreateLogResp(filteredEntries, logOrder)
	return logOutput, nil
}

func filterEntries(entries []types.LogEntry, logDate types.Date, logCategory types.Category) []types.LogEntry {
	filteredEntries := entries
	if logDate.IsSet() {
		filteredEntries = filter(entries, func(entry types.LogEntry) bool {
			startTime := logDate.From()
			endTime := logDate.To()
			if entry.Date.Time.After(startTime) && entry.Date.Time.Before(endTime) || entry.Date.Time.Equal(startTime) || entry.Date.Time.Equal(endTime) {
				return true
			}
			return false
		})
	}
	if logCategory.Value != "" {
		filteredEntries = filter(filteredEntries, func(entry types.LogEntry) bool {
			if entry.Category.Value == logCategory.Value {
				return true
			}
			return false
		})
	}
	return filteredEntries
}

type isRequiredFunc func(types.LogEntry) bool

func filter(entries []types.LogEntry, isRequired isRequiredFunc) []types.LogEntry {
	var filteredEntries = make([]types.LogEntry, 0)
	for _, entry := range entries {
		if isRequired(entry) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries
}
