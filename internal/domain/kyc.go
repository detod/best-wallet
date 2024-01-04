package domain

type KYCStatus string

const (
	KYCStatusPending    KYCStatus = "pending"
	KYCStatusInProgress KYCStatus = "in_progress"
	KYCStatusApproved   KYCStatus = "approved"
	KYCStatusRejected   KYCStatus = "rejected"
)
