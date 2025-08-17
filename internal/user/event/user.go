package event

type DebitEvent struct {
	EventID string  `json:"event_id"`
	UserID  uint    `json:"user_id"`
	Amount  float64 `json:"amount"`
}
