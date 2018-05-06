package cmd

import (
	"acc/db"
	"acc/types"
	"fmt"
)

func PinCmdHandler(pinAmount float64, pinDate types.Date, pinCategory types.Category) (string, error) {
	if err := db.ProcessTransaction(pinAmount, pinDate, pinCategory); err != nil {
		return "", fmt.Errorf("couldn't process transaction: %v", err)
	}
	return BalanceCmdHandler()
}
