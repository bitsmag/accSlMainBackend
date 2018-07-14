package cmd

import (
	"fmt"

	"github.com/bitsmag/accSlMainBackend/src/acc/db"
	"github.com/bitsmag/accSlMainBackend/src/acc/useroutput"
)

func BalanceCmdHandler() (string, error) {
	var balance float64
	var err error
	if balance, err = db.ReadBalance(); err != nil {
		return "", fmt.Errorf("couldn't read balance from database: %v", err)
	}

	balanceString := useroutput.CreateBalanceResp(balance)
	return balanceString, nil
}
