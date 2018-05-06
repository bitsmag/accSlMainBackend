package useroutput

import (
	"acc/stringutil"
	"acc/types"
	"sort"
	"strconv"
	"strings"
)

// CreateBalanceResp returns a string to inform the user about the account balance
func CreateBalanceResp(balance float64) string {
	formattedBalance := strconv.FormatFloat(balance, 'f', 2, 64)
	return "Balance is: " + formattedBalance
}

// CreateLogResp returns a strign containing all relevant logEntries plus additional separation printables (according to types.Order)
func CreateLogResp(entries []types.LogEntry, logOrder types.Order) string {
	if len(entries) == 0 {
		return ""
	}

	var log string

	switch logOrder.Value {
	case "category":
		entries = orderEntriesByCategory(entries)
		groupedLog := separateEntries(entries, func(entry types.LogEntry) string {
			return entry.Category.String()
		})
		log = createLogWithSeparatorLines(groupedLog)
	default: // "year" is default
		entries = orderEntriesByDate(entries)
		groupedLog := separateEntries(entries, func(entry types.LogEntry) string {
			splitDate := strings.Split(entry.Date.String(), ".")
			year := splitDate[len(splitDate)-1]
			return year
		})
		log = createLogWithSeparatorLines(groupedLog)
	}

	return log
}

func orderEntriesByDate(entries []types.LogEntry) []types.LogEntry {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Time.Before(entries[j].Date.Time)
	})
	return entries
}

func orderEntriesByCategory(entries []types.LogEntry) []types.LogEntry {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Category.Value < entries[j].Category.Value
	})
	return entries
}

type getSeparatorFunc func(types.LogEntry) string

func separateEntries(entries []types.LogEntry, getSeparator getSeparatorFunc) map[string][]types.LogEntry {
	var groupedLog = make(map[string][]types.LogEntry)

	currentSepEntries := make([]types.LogEntry, 0)
	currentSepEntries = append(currentSepEntries, entries[0])

	sepOfFirstEntry := getSeparator(entries[0])
	currentSep := sepOfFirstEntry

	for i, entry := range entries {
		if i+1 < len(entries) {
			if getSeparator(entry) == getSeparator(entries[i+1]) {
				currentSepEntries = append(currentSepEntries, entries[i+1])
			} else {
				groupedLog[currentSep] = currentSepEntries
				currentSepEntries = nil
				currentSep = getSeparator(entries[i+1])
				currentSepEntries = append(currentSepEntries, entries[i+1])
			}
		}
	}
	groupedLog[currentSep] = currentSepEntries

	return groupedLog
}

func createLogWithSeparatorLines(groupedLog map[string][]types.LogEntry) string {
	logWithSeparatorLines := ""
	sortedKeys := sortKeys(groupedLog)
	for _, key := range sortedKeys {
		logEntriesForKey := getSeparatorLine(key) + stringifyEntries(groupedLog[key])
		logWithSeparatorLines += logEntriesForKey
	}
	return logWithSeparatorLines
}

func sortKeys(m map[string][]types.LogEntry) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		intI, errI := strconv.Atoi(keys[i])
		intJ, errJ := strconv.Atoi(keys[j])
		if errI != nil && errJ == nil || // i is number and j is text
			errI == nil && errJ == nil && intI < intJ || // both are numbers and i is less than j
			errI != nil && errJ != nil && keys[i] < keys[j] { // both are strings and i is greater than j
			return true
		}
		return false
	})
	return keys
}

func getSeparatorLine(key string) string {
	separatorLine := "======================================" + key + "\n"
	return separatorLine
}

func stringifyEntries(entries []types.LogEntry) string {
	stringEntries := ""
	for _, entry := range entries {
		line := logEntryLine(entry)
		stringEntries += line + "\n"
	}
	return stringEntries
}

func logEntryLine(entry types.LogEntry) string {
	date := formatDate(entry)
	action := formatAction(entry)
	amount := formatAmount(entry)
	line := formatLine(date, action, amount)
	return line
}

func formatDate(entry types.LogEntry) string {
	date := entry.Date.String()
	date = stringutil.RightPad2Len(date, " ", 10)
	return date
}

func formatAction(entry types.LogEntry) string {
	var action string
	if entry.Amount > 0 {
		action = "[pay-in]"
	} else {
		action = "[pay-out]"
	}
	action = stringutil.RightPad2Len(action, " ", 9)
	return action
}

func formatAmount(entry types.LogEntry) string {
	formattedAmount := strconv.FormatFloat(entry.Amount, 'f', 2, 64)
	formattedAmount = stringutil.LeftPad2Len(formattedAmount, " ", 10)
	return formattedAmount
}

func formatLine(date string, action string, amount string) string {
	formattedString := date + "    " + action + "     " + amount
	return formattedString
}
