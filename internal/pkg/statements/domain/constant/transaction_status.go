package constant

type TxStatus string

const (
	SUCCESS TxStatus = "SUCCESS"
	FAILED  TxStatus = "FAILED"
	PENDING TxStatus = "PENDING"
)
