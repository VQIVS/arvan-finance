package rabbit

const (
	// consumers subscribe to these queues
	RefundQueueName = "finance_billing.refund.request"
	DebitQueueName  = "finance_billing.debit.request"

	// producers publish to these queues
	SMSBilledRouting = "billing.debit.completed"
	Exchange         = "amq.topic"
)
