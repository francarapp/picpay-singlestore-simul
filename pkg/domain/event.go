package domain

type Event struct {
	EventName     string
	EventID       string
	CorrelationID string
	UserID        string
	DtCreated     string
	// DtReceived    time.Time
	// DtIngested    time.Time
	Labels  string
	Payload string
}
