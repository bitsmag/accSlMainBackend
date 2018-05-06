package types

import "errors"

// Order is a wrapper for string to implement the pflag.Value interface
type Order struct {
	Value string
}

// String returns the string-value of the category
func (o Order) String() string {
	return o.Value
}

// Set sets the order-value
func (o *Order) Set(order string) error {
	if order == "year" || order == "category" {
		o.Value = order
		return nil
	}
	return errors.New("couldn't parse order - please provide one of the following values ['year' | 'category']")
}

// Type returns the type
func (o Order) Type() string {
	return "order"
}
