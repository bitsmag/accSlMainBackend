package cmd

import (
	"fmt"

	"github.com/bitsmag/accSlMainBackend/src/acc/db"
	"github.com/bitsmag/accSlMainBackend/src/acc/types"
)

func PoutCmdHandler(pinAmount float64, pinDate types.Date, pinCategory types.Category) (string, error) {
	if err := db.ProcessTransaction(pinAmount, pinDate, pinCategory); err != nil {
		return "", fmt.Errorf("couldn't process transaction: %v", err)
	}
	return BalanceCmdHandler()
}
