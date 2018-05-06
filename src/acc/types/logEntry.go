package types

// LogEntry is the type used for logging the bookings
type LogEntry struct {
	Amount   float64
	Date     Date
	Category Category
}
