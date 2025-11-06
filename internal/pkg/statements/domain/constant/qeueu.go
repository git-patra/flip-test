package constant

// ===== Constants (exchange, routing keys, queue names)
const (
	ExchangeTransactions = "transactions"

	RKTransactionsFailed  = "transactions.failed"
	RKTransactionsPending = "transactions.pending"

	QueueReconcileFailed = "reconcile_failed"
	QueueReviewPending   = "review_pending"
)
