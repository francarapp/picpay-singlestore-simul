package domain

type Event struct {
	EventName     string
	EventID       string
	CorrelationID string
	UserID        string
	DtCreated     string
	DtReceived    string
	DtIngested    string
	Labels        string
	Payload       string
}

// TableName overrides the table name used by User to `profiles`
func (Event) TableName() string {
	return "nn_event"
}
