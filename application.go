package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bitsmag/accSlMainBackend/src/acc/cmd"
	"github.com/bitsmag/accSlMainBackend/src/acc/db"
	"github.com/bitsmag/accSlMainBackend/src/acc/types"
)

func main() {
	var err error

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	if arg == "rebuildTables" || arg == "rt" { // rebuild database if "rt" param is present
		if err = db.SetUpTables(); err != nil {
			fmt.Println(err)
			handleError(err)
		}
		if err = db.InsertInitialValues(); err != nil {
			fmt.Println(err)
			handleError(err)
		}
	} else { // check if required "default" balance-entry is available already
		var err error
		if _, err = db.ReadBalance(); err == db.ErrDefaultBalanceDoesNotExist {
			if err = db.InsertInitialValues(); err != nil {
				fmt.Println(err)
				handleError(err)
			}
		}
	}

	http.HandleFunc("/help", help)
	http.HandleFunc("/pin/", pin)
	http.HandleFunc("/pout/", pout)
	http.HandleFunc("/balance", balance)
	http.HandleFunc("/log/", log)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		panic(err)
	}
}

func help(w http.ResponseWriter, r *http.Request) {
	helptext := "help              Prints this helptext." + "\n"
	helptext += "pin <amount>      Pays in a certain amount to the account. FLAGS: -d (dd.mm.yyy); -c <category>" + "\n"
	helptext += "pout <amount>     Pays out a certain amount from the account. FLAGS: -d (dd.mm.yyy); -c <category>" + "\n"
	helptext += "balance           Prints the current balance." + "\n"
	helptext += "log               Prints past transactions. FLAGS: -o ['year' | 'category']; -d (dd.mm.yyyy); -c <category>" + "\n"

	w.Write([]byte(helptext))
}

func pin(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")

	amount, err := strconv.ParseFloat(urlPart[2], 64)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	}

	var date types.Date
	if 3 < len(urlPart) {
		date.Set(urlPart[3])
	}

	var category types.Category
	if 4 < len(urlPart) {
		category.Set(urlPart[4])
	}

	resp, err := cmd.PinCmdHandler(amount, date, category)

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	} else {
		w.Write([]byte(resp))
	}
}

func pout(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")

	amount, err := strconv.ParseFloat(urlPart[2], 64)
	amount = -amount
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	}

	var date types.Date
	if 3 < len(urlPart) {
		date.Set(urlPart[3])
	}

	var category types.Category
	if 4 < len(urlPart) {
		category.Set(urlPart[4])
	}

	resp, err := cmd.PoutCmdHandler(amount, date, category)

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	} else {
		w.Write([]byte(resp))
	}
}

func balance(w http.ResponseWriter, r *http.Request) {
	resp, err := cmd.BalanceCmdHandler()
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	} else {
		w.Write([]byte(resp))
	}
}

func log(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")

	var order types.Order
	if 2 < len(urlPart) {
		order.Set(urlPart[2])
	}

	var date types.Date
	if 3 < len(urlPart) {
		date.Set(urlPart[3])
	}

	var category types.Category
	if 4 < len(urlPart) {
		category.Set(urlPart[4])
	}

	resp, err := cmd.LogCmdHandler(order, date, category)

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Some error occured on the server"))
	} else {
		w.Write([]byte(resp))
	}
}

func handleError(e error) {
	os.Exit(1)
}
