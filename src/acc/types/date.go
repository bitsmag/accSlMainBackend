package types

import (
	"acc/stringutil"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Date is a wrapper for time.Time to implement the pflag.Value interface
type Date struct {
	Time     time.Time
	Accuracy string // day / month / year
}

// String returns the date in the following format DD.MM.YYYY
func (d Date) String() string {
	var day, month, year string
	switch d.Accuracy {
	case "day":
		day = strconv.Itoa(d.Time.Day())
		day = stringutil.LeftPad2Len(day, "0", 2)
		month = strconv.Itoa(int(d.Time.Month()))
		month = stringutil.LeftPad2Len(month, "0", 2)
		year = strconv.Itoa(d.Time.Year())
		year = stringutil.LeftPad2Len(year, "0", 4)
	case "month":
		day = ".."
		month = strconv.Itoa(int(d.Time.Month()))
		month = stringutil.LeftPad2Len(month, "0", 2)
		year = strconv.Itoa(d.Time.Year())
		year = stringutil.LeftPad2Len(year, "0", 4)
	case "year":
		day = ".."
		month = ".."
		year = strconv.Itoa(d.Time.Year())
		year = stringutil.LeftPad2Len(year, "0", 4)
	default:
		return "N/A"
	}
	return day + "." + month + "." + year
}

// Set sets the date and expects the following format DD.MM.YYYY
func (d *Date) Set(date string) error {
	splitDate := strings.FieldsFunc(date, func(r rune) bool {
		return strings.ContainsRune(".", r)
	})

	var day, month, year int
	var err error
	switch len(splitDate) {
	case 3:
		day, err = strconv.Atoi(splitDate[0])
		month, err = strconv.Atoi(splitDate[1])
		year, err = strconv.Atoi(splitDate[2])
		if err != nil || day <= 0 || month <= 0 || year <= 0 {
			return fmt.Errorf("couldn't parse date - please provide numbers in the following format [DD].[MM].YYYY: %v", err)
		}
		d.Accuracy = "day"
	case 2:
		day = 1
		month, err = strconv.Atoi(splitDate[0])
		year, err = strconv.Atoi(splitDate[1])
		if err != nil || month <= 0 || year <= 0 {
			return fmt.Errorf("couldn't parse date - please provide numbers in the following format [DD].[MM].YYYY: %v", err)
		}
		d.Accuracy = "month"
	case 1:
		day = 1
		month = 1
		year, err = strconv.Atoi(splitDate[0])
		if err != nil || year <= 0 {
			return fmt.Errorf("couldn't parse date - please provide numbers in the following format [DD].[MM].YYYY: %v", err)
		}
		d.Accuracy = "year"
	default:
		return errors.New("couldn't parse date - please provide numbers in the following format [DD].[MM].YYYY: %v")
	}

	d.Time = time.Date(year, time.Month(month), day, int(0), int(0), int(0), int(0), time.UTC)
	return nil
}

// Type returns the type
func (d Date) Type() string {
	return "date"
}

// From returns the mimimum time.Time which is within the accuracy of the date
func (d Date) From() time.Time {
	return d.Time
}

// To returns the maximum time.Time which is within the accuracy of the date
func (d Date) To() time.Time {
	var toTime time.Time

	day := 0 // falls back to last day of previous month (hence month needs to be increased by 1)
	month := 13
	year := d.Time.Year()
	if d.Accuracy == "month" {
		month = int(d.Time.Month()) + 1
	}
	if d.Accuracy == "day" {
		month = int(d.Time.Month())
		day = d.Time.Day()
	}

	toTime = time.Date(year, time.Month(month), day, int(23), int(59), int(59), int(999999999), time.UTC)
	return toTime
}

func (d Date) IsSet() bool {
	if d.Accuracy != "" {
		return true
	}
	return false
}
