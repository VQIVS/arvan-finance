package event

type Type string
type Domain string
type Status string

const (
	SMSDebitEvent  Type   = "sms_debit"
	SMSCreditEvent Type   = "sms_credit"
	SMS            Domain = "sms"
	Finance        Domain = "finance"
	StatusSuccess  Status = "success"
	StatusFailed   Status = "failed"
)

type UserBalanceEvent struct {
	Domain Domain  `json:"domain"`
	UserID uint    `json:"user_id"`
	SMSID  uint    `json:"sms_id"`
	Amount float64 `json:"amount"`
	Type   Type    `json:"type"`
}

type SMSUpdateEvent struct {
	Domain Domain `json:"domain"`
	SMSID  uint   `json:"sms_id"`
	Status Status `json:"status"`
}
