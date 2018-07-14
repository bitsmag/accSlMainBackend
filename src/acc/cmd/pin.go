package cmd

import (
	"fmt"

	"github.com/bitsmag/accSlMainBackend/src/acc/types"

	"github.com/bitsmag/accSlMainBackend/src/acc/db"
)

func PinCmdHandler(pinAmount float64, pinDate types.Date, pinCategory types.Category) (string, error) {
	if err := db.ProcessTransaction(pinAmount, pinDate, pinCategory); err != nil {
		return "", fmt.Errorf("couldn't process transaction: %v", err)
	}
	return BalanceCmdHandler()
}
