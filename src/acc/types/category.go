package types

// Category is a wrapper for string to implement the pflag.Value interface
type Category struct {
	Value string
}

// String returns the string-value of the category
func (c Category) String() string {
	if c.Value != "" {
		return c.Value
	}
	return "N/A"

}

// Set sets the category
func (c *Category) Set(category string) error {
	c.Value = category
	return nil
}

// Type returns the type
func (c Category) Type() string {
	return "categroy"
}
