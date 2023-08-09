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
